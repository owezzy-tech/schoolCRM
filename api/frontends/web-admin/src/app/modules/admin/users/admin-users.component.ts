import { AsyncPipe, DatePipe } from '@angular/common';
import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    viewChild,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { FuseDrawerComponent } from '@fuse/components/drawer';
import { toSignal } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { AdminUsersService } from 'app/core/admin-users/admin-users.service';
import {
    AdminUser,
    SCHOOL_ROLES,
    SchoolRole,
    UpdateUserRequest,
} from 'app/core/admin-users/admin-users.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { catchError, of } from 'rxjs';

type DiagramRole = 'Admin' | 'Admissions' | 'Reviewer' | 'Marketing' | 'Read-only';

interface DiagramRoleColumn {
    title: DiagramRole;
    schoolRole: SchoolRole;
}

interface RolePermission {
    capability: string;
    roles: Record<DiagramRole, boolean>;
}

interface EditUserForm {
    name: string;
    email: string;
    department: string;
    enabled: boolean;
    password: string;
    passwordConfirm: string;
}

@Component({
    selector: 'app-admin-users',
    imports: [
        AsyncPipe,
        DatePipe,
        FuseDrawerComponent,
        FormsModule,
        MatButtonModule,
        MatChipsModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatSelectModule,
        MatSlideToggleModule,
    ],
    templateUrl: './admin-users.component.html',
    styleUrl: './admin-users.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AdminUsersComponent {
    private readonly adminUsersService = inject(AdminUsersService);
    private readonly editUserDrawer =
        viewChild.required<FuseDrawerComponent>('editUserDrawer');

    readonly roles = SCHOOL_ROLES;
    readonly diagramRoles: DiagramRoleColumn[] = [
        { title: 'Admin', schoolRole: 'SCHOOL_ADMIN' },
        { title: 'Admissions', schoolRole: 'TEACHER' },
        { title: 'Reviewer', schoolRole: 'TEACHER' },
        { title: 'Marketing', schoolRole: 'TEACHER' },
        { title: 'Read-only', schoolRole: 'PARENT' },
    ];
    readonly users$ = this.adminUsersService.users$;
    readonly usersResult = toSignal(this.users$, {
        initialValue: { items: [], total: 0, page: 1, rowsPerPage: 25 },
    });
    readonly activeUsers = computed(
        () => this.usersResult().items.filter((user) => user.enabled).length
    );
    readonly administratorUsers = computed(
        () =>
            this.usersResult().items.filter((user) =>
                user.roles.some(
                    (role) => role === 'SUPER_ADMIN' || role === 'SCHOOL_ADMIN'
                )
            ).length
    );
    readonly pendingInvites = computed(
        () => this.usersResult().items.filter((user) => !user.enabled).length
    );
    readonly permissions: RolePermission[] = [
        {
            capability: 'View constituents',
            roles: {
                Admin: true,
                Admissions: true,
                Reviewer: true,
                Marketing: true,
                'Read-only': true,
            },
        },
        {
            capability: 'Edit applications',
            roles: {
                Admin: true,
                Admissions: true,
                Reviewer: true,
                Marketing: false,
                'Read-only': false,
            },
        },
        {
            capability: 'Make admission decisions',
            roles: {
                Admin: true,
                Admissions: true,
                Reviewer: true,
                Marketing: false,
                'Read-only': false,
            },
        },
        {
            capability: 'Send campaigns',
            roles: {
                Admin: true,
                Admissions: false,
                Reviewer: false,
                Marketing: true,
                'Read-only': false,
            },
        },
        {
            capability: 'Manage users & roles',
            roles: {
                Admin: true,
                Admissions: false,
                Reviewer: false,
                Marketing: false,
                'Read-only': false,
            },
        },
        {
            capability: 'View audit log',
            roles: {
                Admin: true,
                Admissions: false,
                Reviewer: false,
                Marketing: false,
                'Read-only': false,
            },
        },
    ];

    errorMessage = '';
    editErrorMessage = '';
    editForm: EditUserForm = this.emptyEditForm();
    savingUserID = '';
    search = '';
    selectedUser: AdminUser | null = null;
    selectedRole: SchoolRole | '' = '';

    constructor() {
        this.loadUsers();
    }

    loadUsers(): void {
        this.errorMessage = '';
        this.adminUsersService
            .query({
                page: 1,
                rows: 25,
                orderBy: 'name,ASC',
                name: this.search.trim() || undefined,
            })
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to load users.'
                    );
                    return of(undefined);
                })
            )
            .subscribe();
    }

    inviteUser(): void {
        this.errorMessage =
            'Invite user is a UI placeholder until the backend invite endpoint is available.';
    }

    editUser(user: AdminUser): void {
        this.errorMessage = '';
        this.editErrorMessage = '';
        this.selectedUser = user;
        this.editForm = {
            name: user.name,
            email: user.email,
            department: user.department,
            enabled: user.enabled,
            password: '',
            passwordConfirm: '',
        };
        this.editUserDrawer().open();
    }

    closeEditDrawer(): void {
        this.editUserDrawer().close();
        this.selectedUser = null;
        this.editErrorMessage = '';
        this.editForm = this.emptyEditForm();
    }

    saveEditedUser(): void {
        if (!this.selectedUser) {
            return;
        }

        if (this.editForm.password !== this.editForm.passwordConfirm) {
            this.editErrorMessage = 'Passwords must match before saving.';
            return;
        }

        const request: UpdateUserRequest = {
            name: this.editForm.name.trim(),
            email: this.editForm.email.trim(),
            department: this.editForm.department.trim(),
            enabled: this.editForm.enabled,
        };

        if (this.editForm.password) {
            request.password = this.editForm.password;
            request.passwordConfirm = this.editForm.passwordConfirm;
        }

        this.editErrorMessage = '';
        this.savingUserID = this.selectedUser.id;
        this.adminUsersService
            .update(this.selectedUser.id, request)
            .pipe(
                catchError((error) => {
                    this.editErrorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to update user.'
                    );
                    return of(undefined);
                })
            )
            .subscribe((updatedUser) => {
                this.savingUserID = '';

                if (!updatedUser) {
                    return;
                }

                this.closeEditDrawer();
                this.loadUsers();
            });
    }

    filteredUsers(): AdminUser[] {
        const query = this.search.trim().toLowerCase();

        return this.usersResult().items.filter((user) => {
            const matchesRole = this.selectedRole
                ? user.roles.includes(this.selectedRole as SchoolRole)
                : true;
            const matchesSearch = query
                ? [user.name, user.email, user.department]
                      .join(' ')
                      .toLowerCase()
                      .includes(query)
                : true;

            return matchesRole && matchesSearch;
        });
    }

    initials(user: AdminUser): string {
        return user.name
            .split(' ')
            .filter((part) => part.length > 0)
            .slice(0, 2)
            .map((part) => part[0].toUpperCase())
            .join('');
    }

    roleBadgeClass(role: SchoolRole): string {
        switch (role) {
            case 'SUPER_ADMIN':
                return 'bg-purple-100 text-purple-700';
            case 'SCHOOL_ADMIN':
                return 'bg-indigo-100 text-indigo-700';
            case 'TEACHER':
                return 'bg-cyan-100 text-cyan-700';
            case 'STUDENT':
                return 'bg-orange-100 text-orange-700';
            case 'PARENT':
                return 'bg-slate-100 text-slate-700';
        }
    }

    private emptyEditForm(): EditUserForm {
        return {
            name: '',
            email: '',
            department: '',
            enabled: true,
            password: '',
            passwordConfirm: '',
        };
    }

    updateRoles(user: AdminUser, roles: SchoolRole[]): void {
        this.errorMessage = '';
        this.savingUserID = user.id;
        this.adminUsersService
            .updateRoles(user.id, { roles })
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to update roles.'
                    );
                    return of(undefined);
                })
            )
            .subscribe(() => {
                this.savingUserID = '';
                this.loadUsers();
            });
    }

    deleteUser(user: AdminUser): void {
        this.errorMessage = '';
        this.savingUserID = user.id;
        this.adminUsersService
            .delete(user.id)
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to remove user'
                    );
                    return of(undefined);
                })
            )
            .subscribe(() => {
                this.savingUserID = '';
                this.loadUsers();
            });
    }
}
