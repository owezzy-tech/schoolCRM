// Opt-in public-HTTPS acceptance using separate synthetic principals.
// Credentials are injected through environment variables and never printed.
import { createHash, randomUUID } from 'node:crypto';
import assert from 'node:assert/strict';

const base = 'https://web-admin-staging-c78a.up.railway.app';
const school = 'ed78d739-d93a-431d-8957-da3e99c78cef';
const department = '495cf156-ee17-4d06-b271-c4ccbaa47fe4';
const plan = '495b0af4-b207-4cf9-b321-382e9ca5cc00';
const source = 'b1f3d4e7-1c1d-4811-a18e-b87a5ba6a7c5';
const generation = 'b99268af-c194-47c1-80aa-07b731d2b815';
const reuseRequest = '0a5c6da0-07a4-4d2d-a4eb-0a8ea8169c77';
const resource = value => ({ ...value.attributes, id: value.id });

async function request(method, path, token, body, expected = [200, 201, 204]) {
    const response = await fetch(base + path, {
        method,
        headers: { ...(token ? { Authorization: 'Bearer ' + token } : {}),
            ...(body ? { 'Content-Type': 'application/json' } : {}) },
        ...(body ? { body: JSON.stringify(body) } : {}),
        signal: AbortSignal.timeout(30000),
    });
    if (!expected.includes(response.status)) {
        throw new Error(`Acceptance failed at ${method} ${path}: HTTP ${response.status}`);
    }
    return response.status === 204 ? null : response.json();
}

async function login(email, password) {
    const data = await request('POST', '/v1/auth/login', null, { email, password });
    return data.data.attributes;
}

async function verify() {
    if (process.env.SCHOOLCRM_HOSTED_PROOF !== 'staging') throw new Error('Staging opt-in required');
    const password = process.env.SCHOOLCRM_STAGING_QA_PASSWORD;
    if (!password || !process.env.SCHOOLCRM_STAGING_ADMIN_PASSWORD) throw new Error('Private credentials required');
    const manager = await login('owezzy@owenadirah.com', process.env.SCHOOLCRM_STAGING_ADMIN_PASSWORD);
    const admin = manager.accessToken;
    const identities = (await request('GET', '/v1/users?rows=100', admin)).data.map(resource);
    async function principal(label) {
        const email = `schoolai-${label}@staging.invalid`;
        let user = identities.find(u => u.email === email);
        if (!user) {
            user = resource((await request('POST', '/v1/users', admin, {
                name: 'Staging AI ' + label, email, roles: ['TEACHER'], department: 'Languages',
                password, passwordConfirm: password,
            })).data);
        }
        return login(email, password);
    }
    const teacher = await principal('teacher');
    const hod = await principal('hod');
    const dean = await principal('dean');
    const outsider = await principal('outsider');
    assert.equal(new Set([teacher.user.id, hod.user.id, dean.user.id, outsider.user.id]).size, 4);
    const memberships = (await request('GET', `/v1/schools/${school}/memberships`, admin)).data.map(resource);
    for (const [actor, capability] of [[hod, 'teach'], [hod, 'review_lessons'], [dean, 'approve_lessons']]) {
        if (!memberships.some(m => m.userID === actor.user.id && m.departmentID === department && m.capability === capability && m.active)) {
            await request('POST', `/v1/schools/${school}/memberships`, admin, {
                userID: actor.user.id, departmentID: department, capability,
            });
        }
    }
    await request('GET', `/v1/rag/curriculum/sources/${source}`, outsider.accessToken, null, [403, 404]);
    await request('GET', `/v1/rag/lessons/generations/${generation}`, outsider.accessToken, null, [404]);
    await request('GET', `/v1/rag/curriculum/sources/${source}/original`, null, null, [401]);
    const original = await fetch(`${base}/v1/rag/curriculum/sources/${source}/original`, {
        headers: { Authorization: 'Bearer ' + teacher.accessToken }, signal: AbortSignal.timeout(30000),
    });
    assert.equal(original.status, 200);
    assert.equal(createHash('sha256').update(Buffer.from(await original.arrayBuffer())).digest('hex'),
        'eeea001277e1d952f4a6ebd4e42a21072cbb1e261d2e7228e50fe699ef041d0d');
    // Every run gets a new probe copied from the real published model output.
    // No model call is needed; all denial checks and transitions execute afresh.
    const probe = resource((await request('POST', `/v1/lessons/${plan}/reuse`, teacher.accessToken, {
        version: 1, requestID: randomUUID(),
    })).data);
    assert.equal(probe.status, 'draft');
    console.log(JSON.stringify({ phase: 'fresh_workflow_probe', probe: probe.id }));
    const versionsPath = `/v1/lessons/${probe.id}/versions`;
    const denials = [];
    async function denied(name, method, path, actor, body, statuses) {
        await request(method, path, actor.accessToken, body, statuses);
        denials.push(name);
    }
    await denied('private_read', 'GET', `/v1/lessons/${probe.id}`, hod, null, [404]);
    await denied('private_reuse', 'POST', `/v1/lessons/${probe.id}/reuse`, hod,
        { version: 1, requestID: randomUUID() }, [404]);
    await denied('premature_publication', 'POST', versionsPath + '/1/publish', teacher, null, [409]);
    await request('POST', versionsPath + '/1/submit', teacher.accessToken);
    await denied('teacher_review', 'POST', versionsPath + '/1/review', teacher,
        { decision: 'approve', feedback: 'Not authorised' }, [403]);
    await request('POST', versionsPath + '/1/review', hod.accessToken,
        { decision: 'approve', feedback: 'Hosted reviewed evidence and lesson' });
    await denied('hod_dean_approval', 'POST', versionsPath + '/1/approval', hod,
        { decision: 'approve', feedback: 'Not authorised' }, [403]);
    await request('POST', versionsPath + '/1/approval', dean.accessToken,
        { decision: 'approve', feedback: 'Hosted approval by distinct dean' });
    await request('POST', versionsPath + '/1/publish', teacher.accessToken);
    const reused = resource((await request('POST', `/v1/lessons/${plan}/reuse`, hod.accessToken, {
        version: 1, requestID: reuseRequest,
    })).data);
    assert.equal(reused.status, 'draft');
    assert.equal(reused.publishedVersion, null);
    await request('GET', `/v1/lessons/${reused.id}`, teacher.accessToken, null, [404]);
    let versions = (await request('GET', versionsPath, teacher.accessToken)).data.map(resource);
    const publishedSnapshot = versions.find(v => v.version === 1);
    assert.equal(versions.length, 1);
    await request('POST', versionsPath, teacher.accessToken, {
        baseVersion: 1, title: 'Teacher-edited staging lesson', changeSummary: 'Hosted human edit',
        content: { ...publishedSnapshot.content, objectives: ['Teacher-edited staging objective'] },
    });
    versions = (await request('GET', versionsPath, teacher.accessToken)).data.map(resource);
    const old = versions.find(v => v.version === 1);
    assert.equal(old.status, 'published');
    assert.equal(old.reviewerID, hod.user.id);
    assert.equal(old.approverID, dean.user.id);
    assert.equal(old.authorID, teacher.user.id);
    assert.equal(versions.find(v => v.version === 2).status, 'draft');
    assert.notDeepEqual(old.content.objectives, versions.find(v => v.version === 2).content.objectives);
    const updated = resource((await request('GET', `/v1/lessons/${probe.id}`, teacher.accessToken)).data);
    assert.equal(updated.currentVersion, 2);
    assert.equal(updated.publishedVersion, 1);
    assert.equal(denials.length, 5);
    console.log(JSON.stringify({ phase: 'workflow_complete', plan, probe: probe.id, reused: reused.id, teacher: teacher.user.id,
        hod: hod.user.id, dean: dean.user.id, currentVersion: 2, publishedVersion: 1,
        scopeDenials: 'passed', observedWorkflowDenials: denials, sourceChecksum: 'verified' }));
}

verify().catch(error => {
    // Assertion diffs and transport errors could contain response data. Suppress them.
    console.error(error.message.startsWith('Acceptance failed') ? error.message : 'Staging acceptance assertion or transport failed');
    process.exitCode = 1;
});
