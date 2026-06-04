import { expect, Page, test } from '@playwright/test';

const staffUser = {
    id: '5cf37266-3473-4006-984f-9325122678b7',
    name: 'Admin Gopher',
    email: 'admin@example.com',
    roles: ['SCHOOL_ADMIN', 'ADMISSIONS_ADMIN', 'APPLICATION_REVIEWER'],
};

const applicantFixture = {
    applicationID: 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
    constituentID: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    programID: 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    academicTermID: 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
};

test.describe('Admissions staff experience', () => {
    test.beforeEach(async ({ page }) => {
        await mockAuth(page);
        await signIn(page);
    });

    test('opens the applications list and application detail happy path', async ({
        page,
    }) => {
        const requests = await mockAdminApplications(page);

        await page.goto('/applications');

        await expect(
            page.getByRole('heading', { name: 'Applications' })
        ).toBeVisible();
        await expect(
            page.getByPlaceholder('Search applications...')
        ).toBeVisible();
        await expect(
            page.getByRole('button', { name: /New application/i })
        ).toBeVisible();
        await expect(page.getByText('Bachelor of Commerce')).toBeVisible();
        await expect(page.getByText('BCOM-QA')).toBeVisible();
        await expect(page.getByText('2026 Main Intake')).toBeVisible();
        await expect(page.getByText('IN REVIEW')).toBeVisible();

        await page
            .getByRole('link', { name: applicantFixture.applicationID })
            .click();

        await expect(page).toHaveURL(
            new RegExp(`/applications/${applicantFixture.applicationID}$`)
        );
        await expect(
            page.getByRole('link', { name: /Back to applications/i })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: applicantFixture.constituentID })
        ).toBeVisible();
        await expect(
            page.getByText('Application', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText('Documents', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'KUCCPS placement' })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'KCSE result snapshot' })
        ).toBeVisible();
        await page.getByRole('tab', { name: 'Timeline' }).click();
        await expect(
            page.getByText('READY FOR REVIEW → IN REVIEW')
        ).toBeVisible();
        await page.getByRole('tab', { name: 'Notes' }).click();
        await expect(
            page.getByText('Reviewer accepted the application for evaluation.')
        ).toBeVisible();
        expect(requests.loadedApplications).toBe(1);
        expect(requests.loadedApplicationDetail).toBe(1);
        expect(requests.loadedTransitions).toBe(1);
    });

    test('shows the staff review workspace with authorized application queue', async ({
        page,
    }) => {
        await page.goto('/applications/review');

        await expect(
            page.getByText('Applications / Review Workspace')
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Staff review workspace' })
        ).toBeVisible();
        await expect(page.getByText('Authorized queue')).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Assigned applications' })
        ).toBeVisible();
        await expect(
            page.getByRole('button', { name: /Open application APP-3024/i })
        ).toBeVisible();
        await expect(page.getByText('Decision due')).toBeVisible();
    });

    test('navigates from lead pipeline to admissions settings', async ({
        page,
    }) => {
        await page.goto('/leads');

        await expect(page.getByText('Admissions / Leads')).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Lead pipeline' })
        ).toBeVisible();
        await expect(
            page.getByRole('combobox', { name: 'Score band' })
        ).toBeVisible();

        await page
            .getByRole('button', { name: 'Configure scoring rules' })
            .click();

        await expect(page).toHaveURL(/\/admissions-settings$/);
        await expect(page.getByText('Admissions / Settings')).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Settings & configuration' })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Custom field registry' })
        ).toBeVisible();
        await expect(
            page.getByText('Lead Scoring Rules', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText('Import/Export', { exact: true })
        ).toBeVisible();
    });
});

test.describe('Applicant portal experience', () => {
    test('starts the public application happy path', async ({ page }) => {
        const applicantRequests = await mockApplicantAdmissions(page);

        await page.goto('/portal');

        await expect(
            page.getByRole('link', { name: /SchoolCRM Kenya Applicant/i })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', {
                name: 'Start your admissions journey to School Of Thought University.',
            })
        ).toBeVisible();
        await expect(
            page.getByText(/KUCCPS placement or self-sponsored intake/i)
        ).toBeVisible();
        await expect(
            page.getByRole('link', { name: /Start application/i })
        ).toBeVisible();
        await expect(page.getByText('Email-only portal access')).toBeVisible();

        await page.getByRole('link', { name: /Start application/i }).click();

        await expect(page).toHaveURL(/\/portal\/apply$/);
        await expect(
            page.getByRole('heading', {
                name: 'Apply for the 2026 main intake',
            })
        ).toBeVisible();
        await expect(
            page.getByLabel('Application progress: step 1 of 6, account')
        ).toBeVisible();
        await expect(
            page.getByText('Applicant account', { exact: true })
        ).toBeVisible();
        await expect(page.getByLabel('Preferred programme')).not.toBeVisible();

        await page.getByLabel('First name').fill('John');
        await page.getByLabel('Last name').fill('Applicant');
        await page.getByLabel('Email').fill('applicant@example.com');
        await page.getByLabel('Phone').fill('+254712345678');
        await page.getByLabel('Date of birth').fill('2005-01-15');
        await page.getByLabel('Password', { exact: true }).fill('gophers');
        await page.getByLabel('Confirm password').fill('gophers');

        await expect(
            page.getByRole('button', { name: /Continue/i })
        ).toBeVisible();
        await expect(page.getByText('Application fee')).toBeVisible();
        await expect(page.getByText(/Ksh\s*150\.00/i)).toBeVisible();
        await expect(page.getByText(/Payment channel: M-Pesa/i)).toBeVisible();
        await expect(page.getByText('Need help?')).toBeVisible();
        await expect(
            page.getByText(/admissions@schoolcrm.ac.ke/i)
        ).toBeVisible();

        await page.getByRole('button', { name: /Continue/i }).click();
        await expect(
            page.getByText('Applicant account created.')
        ).toBeVisible();
        await expect(
            page.getByLabel('Application progress: step 2 of 6, programme')
        ).toBeVisible();
        await expect(
            page.getByText('Programme selection', { exact: true })
        ).toBeVisible();
        await expect(page.getByLabel('Preferred programme')).toBeVisible();
        await expect(page.getByLabel('Sponsorship route')).toBeVisible();
        await expect(page.getByLabel('Preferred programme')).toHaveValue(
            'Bachelor of Commerce'
        );
        await expect(page.getByLabel('Intake')).toHaveValue('2026 Main Intake');

        await page.getByRole('button', { name: /Continue/i }).click();
        await expect(
            page.getByText('Draft saved to admissions.')
        ).toBeVisible();
        await expect(
            page.getByLabel('Application progress: step 3 of 6, kcse details')
        ).toBeVisible();
        await expect(page.getByLabel('KCSE index number')).toBeVisible();
        await expect(page.getByLabel('KCSE mean grade')).toBeVisible();
        await page.getByLabel('KCSE index number').fill('12345678901/2025');
        await page.getByLabel('KCSE year').fill('2025');
        await page.getByLabel('KCSE mean grade').fill('B+');
        await page
            .getByLabel('Subject highlights')
            .fill('Maths B+, English A-');

        await page.getByRole('button', { name: /Documents/i }).click();
        await expect(
            page.getByLabel('Application progress: step 5 of 6, documents')
        ).toBeVisible();
        await expect(
            page.getByText('Supporting documents', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText('KCSE result document', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText('M-Pesa receipt', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText(/Drop KCSE result document or/i)
        ).toBeVisible();
        await expect(page.getByText(/Drop M-Pesa receipt or/i)).toBeVisible();

        await page.getByRole('button', { name: /Back/i }).click();
        await expect(
            page.getByLabel('Application progress: step 4 of 6, placement')
        ).toBeVisible();
        await expect(
            page.getByText('KUCCPS or self-sponsored placement', {
                exact: true,
            })
        ).toBeVisible();
        await expect(page.getByLabel('KUCCPS placement number')).toBeVisible();
        await expect(page.getByLabel('National ID or passport')).toBeVisible();
        await page
            .getByLabel('KUCCPS placement number')
            .fill('KUCCPS-2026-0001');
        await page.getByLabel('National ID or passport').fill('12345678');

        await page.getByRole('button', { name: /Review/i }).click();
        await expect(
            page.getByLabel('Application progress: step 6 of 6, review')
        ).toBeVisible();
        await expect(page.getByLabel('M-Pesa confirmation code')).toBeVisible();
        await expect(page.getByLabel('Review notes')).toBeVisible();
        await page.getByLabel('M-Pesa confirmation code').fill('QF123ABC45');
        await page.getByLabel('Review notes').fill('Ready for review');
        await page.getByRole('button', { name: /Submit application/i }).click();

        await expect(
            page.getByText('Application submitted for admissions review.')
        ).toBeVisible();
        expect(applicantRequests.onboardedApplicants).toBe(1);
        expect(
            applicantRequests.createdApplications +
                applicantRequests.savedApplications
        ).toBeGreaterThanOrEqual(1);
        expect(applicantRequests.submittedApplications).toBe(1);
    });

    test('shows Kenya-localized application status details', async ({
        page,
    }) => {
        await page.goto('/portal/status');

        await expect(
            page.getByRole('link', { name: /SchoolCRM Kenya Applicant/i })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Application status' })
        ).toBeVisible();
        await expect(page.getByText(/2026 Main Intake/i)).toBeVisible();
        await expect(
            page.getByText(/M-Pesa processing is not active in this preview/i)
        ).toBeVisible();
        await expect(page.getByText(/Ksh\s*150\.00/i)).toBeVisible();
        await expect(page.getByText('M-Pesa', { exact: true })).toBeVisible();
        await expect(
            page.getByText('KCSE result slip', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText(
                /Official KCSE result slip or certificate from KNEC/i
            )
        ).toBeVisible();
        await expect(page.getByText('KUCCPS placement letter')).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Reach your admissions officer' })
        ).toBeVisible();
        await expect(
            page.getByText(/admissions@schoolcrm.ac.ke/i)
        ).toBeVisible();
    });
});

async function mockAdminApplications(page: Page): Promise<{
    loadedApplications: number;
    loadedApplicationDetail: number;
    loadedTransitions: number;
}> {
    const requests = {
        loadedApplications: 0,
        loadedApplicationDetail: 0,
        loadedTransitions: 0,
    };

    await page.route('**/v1/admissions/programs*', async (route) => {
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: [
                    {
                        id: applicantFixture.programID,
                        type: 'programs',
                        attributes: {
                            externalSISID: 'KUCCPS-BCOM-2026-QA',
                            name: 'Bachelor of Commerce',
                            code: 'BCOM-QA',
                            degreeLevel: 'BACHELOR',
                            active: true,
                            dateCreated: new Date().toISOString(),
                            dateUpdated: new Date().toISOString(),
                        },
                    },
                ],
                meta: { total: 1, page: 1, rowsPerPage: 100 },
            }),
        });
    });

    await page.route('**/v1/admissions/academic-terms*', async (route) => {
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: [
                    {
                        id: applicantFixture.academicTermID,
                        type: 'academic-terms',
                        attributes: {
                            externalSISID: 'INTAKE-2026-QA',
                            name: '2026 Main Intake',
                            code: 'INT2026-QA',
                            termType: 'MAIN',
                            startDate: '2026-01-01T00:00:00Z',
                            endDate: '2026-12-31T00:00:00Z',
                            active: true,
                            dateCreated: new Date().toISOString(),
                            dateUpdated: new Date().toISOString(),
                        },
                    },
                ],
                meta: { total: 1, page: 1, rowsPerPage: 100 },
            }),
        });
    });

    await page.route(
        `**/v1/admissions/applications/${applicantFixture.applicationID}/transitions*`,
        async (route) => {
            requests.loadedTransitions += 1;
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: [
                        {
                            id: 'transition-in-review',
                            type: 'application-transitions',
                            attributes: {
                                applicationID: applicantFixture.applicationID,
                                fromStatus: 'READY_FOR_REVIEW',
                                toStatus: 'IN_REVIEW',
                                actorID: staffUser.id,
                                reason: 'Reviewer accepted the application for evaluation.',
                                note: 'Reviewer accepted the application for evaluation.',
                                metadata: {},
                                dateCreated: new Date().toISOString(),
                            },
                        },
                    ],
                    meta: { total: 1, page: 1, rowsPerPage: 100 },
                }),
            });
        }
    );

    await page.route(
        `**/v1/admissions/applications/${applicantFixture.applicationID}`,
        async (route) => {
            requests.loadedApplicationDetail += 1;
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: applicantApplicationBody('IN_REVIEW'),
            });
        }
    );

    await page.route('**/v1/admissions/applications*', async (route) => {
        requests.loadedApplications += 1;
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: [
                    {
                        id: applicantFixture.applicationID,
                        type: 'applications',
                        attributes: {
                            constituentID: applicantFixture.constituentID,
                            programID: applicantFixture.programID,
                            academicTermID: applicantFixture.academicTermID,
                            applicationType: 'KUCCPS_PLACEMENT',
                            status: 'IN_REVIEW',
                            assignedReviewerID: staffUser.id,
                            kuccpsPlacement: {
                                placementID: 'KUCCPS-2026-0001',
                                institutionCode: 'STU',
                                programmeCode: 'BCOM-QA',
                                programmeName: 'Bachelor of Commerce',
                                placementYear: 2026,
                            },
                            kcseResult: {
                                indexNumber: '12345678901/2025',
                                examYear: 2025,
                                subjects: [
                                    {
                                        subjectCode: 'MATHS',
                                        grade: 'B+',
                                        points: 10,
                                    },
                                    {
                                        subjectCode: 'ENGLISH',
                                        grade: 'A-',
                                        points: 11,
                                    },
                                ],
                                meanGrade: 'B+',
                                meanPoints: 70,
                            },
                            submittedAt: new Date().toISOString(),
                            dateCreated: new Date().toISOString(),
                            dateUpdated: new Date().toISOString(),
                        },
                    },
                ],
                meta: { total: 1, page: 1, rowsPerPage: 50 },
            }),
        });
    });

    return requests;
}

async function mockApplicantAdmissions(page: Page): Promise<{
    onboardedApplicants: number;
    createdApplications: number;
    savedApplications: number;
    submittedApplications: number;
}> {
    const requests = {
        onboardedApplicants: 0,
        createdApplications: 0,
        savedApplications: 0,
        submittedApplications: 0,
    };
    const applicantToken = createJwt({
        sub: 'f47ac10b-58cc-4372-a567-0e02b2c3d479',
        roles: ['APPLICANT'],
        portal: {
            scope: 'applicant_portal',
            applicationID: '',
            constituentID: applicantFixture.constituentID,
            email: 'applicant@example.com',
        },
    });

    await page.route('**/v1/auth/applicant-portal/onboard', async (route) => {
        requests.onboardedApplicants += 1;
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: {
                    id: applicantFixture.constituentID,
                    type: 'applicant-portal-token',
                    attributes: {
                        accessToken: applicantToken,
                        tokenType: 'Bearer',
                        expiresAt: new Date(
                            Date.now() + 3_600_000
                        ).toISOString(),
                        expiresIn: 3600,
                        applicationID: '',
                        constituentID: applicantFixture.constituentID,
                        applicantName: 'John Applicant',
                        email: 'applicant@example.com',
                    },
                },
            }),
        });
    });

    await page.route('**/v1/admissions/applicant/programs*', async (route) => {
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: [
                    {
                        id: applicantFixture.programID,
                        type: 'programs',
                        attributes: {
                            externalSISID: 'SIS-BCOM',
                            name: 'Bachelor of Commerce',
                            code: 'BCOM',
                            degreeLevel: 'UNDERGRADUATE',
                            active: true,
                            dateCreated: new Date().toISOString(),
                            dateUpdated: new Date().toISOString(),
                        },
                    },
                ],
                meta: { total: 1, page: 1, rowsPerPage: 100 },
            }),
        });
    });

    await page.route(
        '**/v1/admissions/applicant/academic-terms*',
        async (route) => {
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: [
                        {
                            id: applicantFixture.academicTermID,
                            type: 'academic-terms',
                            attributes: {
                                externalSISID: 'SIS-2026-MAIN',
                                name: '2026 Main Intake',
                                code: '2026-MAIN',
                                startDate: '2026-09-01T00:00:00Z',
                                endDate: '2027-04-30T00:00:00Z',
                                active: true,
                                dateCreated: new Date().toISOString(),
                                dateUpdated: new Date().toISOString(),
                            },
                        },
                    ],
                    meta: { total: 1, page: 1, rowsPerPage: 100 },
                }),
            });
        }
    );

    await page.route(
        '**/v1/admissions/applicant/application-form-templates*',
        async (route) => {
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: [
                        {
                            id: 'template-2026-main',
                            type: 'application-form-templates',
                            attributes: {
                                programID: applicantFixture.programID,
                                academicTermID: applicantFixture.academicTermID,
                                applicationType: 'KUCCPS_PLACEMENT',
                                name: '2026 Main Intake Form',
                                version: 1,
                                requiredFields: [],
                                checklistItems: [],
                                active: true,
                                priority: 1,
                                dateCreated: new Date().toISOString(),
                                dateUpdated: new Date().toISOString(),
                            },
                        },
                    ],
                    meta: { total: 1, page: 1, rowsPerPage: 100 },
                }),
            });
        }
    );

    await page.route(
        '**/v1/admissions/applicant/applications',
        async (route) => {
            if (route.request().method() === 'POST') {
                requests.createdApplications += 1;
                await route.fulfill({
                    contentType: 'application/vnd.api+json',
                    body: applicantApplicationBody('DRAFT'),
                });
                return;
            }

            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: [],
                    meta: { total: 0, page: 1, rowsPerPage: 100 },
                }),
            });
        }
    );

    await page.route(
        '**/v1/admissions/applicant/applications/*/transitions',
        async (route) => {
            requests.submittedApplications += 1;
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: applicantApplicationBody('SUBMITTED'),
            });
        }
    );

    await page.route(
        '**/v1/admissions/applicant/applications/*',
        async (route) => {
            if (route.request().method() === 'PUT') {
                requests.savedApplications += 1;
            }

            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: applicantApplicationBody('DRAFT'),
            });
        }
    );

    return requests;
}

async function mockAuth(page: Page): Promise<void> {
    const token = createJwt({ sub: staffUser.id, roles: staffUser.roles });

    await page.route('**/v1/admissions/**', async (route) => {
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: [],
                meta: { total: 0, page: 1, rows: 50 },
            }),
        });
    });

    await page.route('**/v1/auth/login', async (route) => {
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: {
                    type: 'auth-token',
                    attributes: {
                        accessToken: token,
                        tokenType: 'Bearer',
                        expiresAt: new Date(
                            Date.now() + 3_600_000
                        ).toISOString(),
                        expiresIn: 3600,
                        user: staffUser,
                    },
                },
            }),
        });
    });

    await page.route('**/v1/auth/authenticate', async (route) => {
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: {
                    type: 'auth-session',
                    attributes: {
                        userID: staffUser.id,
                        claims: { roles: staffUser.roles },
                        user: staffUser,
                    },
                },
            }),
        });
    });
}

async function signIn(page: Page): Promise<void> {
    await page.goto('/sign-in');
    await page.getByLabel('Email address').fill(staffUser.email);
    await page.getByLabel('Password').fill('gophers');
    await page.getByRole('button', { name: 'Sign in' }).click();
    await expect(page).toHaveURL(/\/dashboard$/);
}

function applicantApplicationBody(
    status: 'DRAFT' | 'SUBMITTED' | 'IN_REVIEW'
): string {
    return JSON.stringify({
        jsonapi: { version: '1.1' },
        data: {
            id: applicantFixture.applicationID,
            type: 'applications',
            attributes: {
                constituentID: applicantFixture.constituentID,
                programID: applicantFixture.programID,
                academicTermID: applicantFixture.academicTermID,
                applicationType: 'KUCCPS_PLACEMENT',
                status,
                kuccpsPlacement: {
                    placementID: 'KUCCPS-2026-0001',
                    institutionCode: 'STU',
                    programmeCode: 'BCOM',
                    programmeName: 'Bachelor of Commerce',
                    placementYear: 2025,
                },
                kcseResult: {
                    indexNumber: '12345678901/2025',
                    examYear: 2025,
                    subjects: [
                        { subjectCode: 'MATHS', grade: 'B+', points: 10 },
                        { subjectCode: 'ENGLISH', grade: 'A-', points: 11 },
                    ],
                    meanGrade: 'B+',
                    meanPoints: 10,
                },
                submittedAt:
                    status === 'DRAFT' ? undefined : new Date().toISOString(),
                dateCreated: new Date().toISOString(),
                dateUpdated: new Date().toISOString(),
            },
        },
    });
}

function createJwt(payload: Record<string, unknown>): string {
    const now = Math.floor(Date.now() / 1000);
    const header = { alg: 'none', typ: 'JWT' };
    const claims = {
        iss: 'playwright',
        iat: now,
        exp: now + 3600,
        ...payload,
    };

    return `${base64Url(header)}.${base64Url(claims)}.`;
}

function base64Url(value: Record<string, unknown>): string {
    return Buffer.from(JSON.stringify(value))
        .toString('base64')
        .replace(/=/g, '')
        .replace(/\+/g, '-')
        .replace(/\//g, '_');
}
