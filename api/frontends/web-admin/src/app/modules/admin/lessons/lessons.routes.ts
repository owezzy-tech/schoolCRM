import { Routes } from '@angular/router';

import { LessonEditorComponent } from './editor/lesson-editor.component';
import { LessonsComponent } from './lessons.component';
import { LessonPlanComponent } from './plan/lesson-plan.component';
import { LessonVersionComponent } from './version/lesson-version.component';

export default [
    { path: '', component: LessonsComponent },
    { path: 'new', component: LessonEditorComponent },
    { path: ':planId', component: LessonPlanComponent },
    { path: ':planId/edit', component: LessonEditorComponent },
    { path: ':planId/versions/:version', component: LessonVersionComponent },
] as Routes;
