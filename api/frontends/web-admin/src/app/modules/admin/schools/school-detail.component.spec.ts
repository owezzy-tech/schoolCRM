import { TestBed } from '@angular/core/testing';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { AdminUsersService } from 'app/core/admin-users/admin-users.service';
import { SchoolAccessService } from 'app/core/school-access/school-access.service';
import { UserService } from 'app/core/user/user.service';
import { of } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { SchoolDetailComponent } from './school-detail.component';

describe('SchoolDetailComponent', () => {
    let grant: ReturnType<typeof vi.fn>;
    let memberships: ReturnType<typeof vi.fn>;

    async function create(roles: string[], ownCapability?: string) {
        grant = vi.fn(() => of({}));
        memberships = vi.fn(() => of([]));
        await TestBed.configureTestingModule({
            imports: [SchoolDetailComponent],
            providers: [
                {
                    provide: SchoolAccessService,
                    useValue: {
                        schools: () => of([{ id: 'school-1', name: 'Nairobi Hill School', dateCreated: '' }]),
                        departments: () => of([{ id: 'sciences', schoolID: 'school-1', name: 'Sciences', dateCreated: '' }]),
                        ownMemberships: () =>
                            of(
                                ownCapability
                                    ? [{ id: 'm', schoolID: 'school-1', userID: 'me', capability: ownCapability, active: true, dateUpdated: '' }]
                                    : []
                            ),
                        memberships,
                        grant,
                    },
                },
                { provide: AdminUsersService, useValue: { query: () => of({ items: [], total: 0, page: 1, rowsPerPage: 100 }) } },
                { provide: FuseConfirmationService, useValue: { open: vi.fn() } },
                { provide: MatSnackBar, useValue: { open: vi.fn() } },
                { provide: UserService, useValue: { user$: of({ id: 'me', name: 'Me', email: 'me@x', roles }) } },
                {
                    provide: ActivatedRoute,
                    useValue: { snapshot: { paramMap: convertToParamMap({ schoolId: 'school-1' }) } },
                },
            ],
        })
            .overrideComponent(SchoolDetailComponent, { set: { template: '' } })
            .compileComponents();
        return TestBed.createComponent(SchoolDetailComponent).componentInstance;
    }

    it('hides membership administration from members without delegation', async () => {
        const component = await create(['SCHOOL_ADMIN']);

        expect(component.canManage()).toBe(false);
        expect(memberships).not.toHaveBeenCalled();
    });

    it('lets a delegated school admin grant department capabilities', async () => {
        const component = await create(['SCHOOL_ADMIN'], 'manage_members');
        component.updateGrant({ userID: 'grace', capability: 'teach' });

        expect(component.canManage()).toBe(true);
        expect(component.grantReady()).toBe(false);

        component.updateGrant({ departmentID: 'sciences' });
        component.submitGrant();

        expect(grant).toHaveBeenCalledWith('school-1', {
            userID: 'grace',
            capability: 'teach',
            departmentID: 'sciences',
        });
    });

    it('sends management delegation without a department', async () => {
        const component = await create(['SUPER_ADMIN']);
        component.updateGrant({ userID: 'admin', capability: 'teach', departmentID: 'sciences' });
        component.updateGrant({ capability: 'manage_members' });
        component.submitGrant();

        expect(grant).toHaveBeenCalledWith('school-1', { userID: 'admin', capability: 'manage_members' });
    });
});
