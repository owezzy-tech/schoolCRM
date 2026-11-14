// Point every SchoolCRM image at the version in package.json. Runs from the
// `version` npm lifecycle, so `npm version <x>` bumps both in one commit.
import { readFileSync, writeFileSync } from 'node:fs';

const { version } = JSON.parse(readFileSync('package.json', 'utf8'));
const file = 'zarf/k8s/vps/release/kustomization.yaml';
const tagged = readFileSync(file, 'utf8').replace(/newTag: .*/g, `newTag: v${version}`);
writeFileSync(file, tagged);
console.log(`${file}: images tagged v${version}`);
