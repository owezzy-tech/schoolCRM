import { Routes } from '@angular/router';

import { SchoolDetailComponent } from './school-detail.component';
import { SchoolsComponent } from './schools.component';

export default [
    { path: '', component: SchoolsComponent },
    { path: ':schoolId', component: SchoolDetailComponent },
] as Routes;
