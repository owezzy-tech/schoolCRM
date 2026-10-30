import { DatePipe } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTabsModule } from '@angular/material/tabs';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { AdminUsersService } from 'app/core/admin-users/admin-users.service';
import { AdminUser } from 'app/core/admin-users/admin-users.types';
import { jsonApiErrorMessage, JsonApiErrorResponse } from 'app/core/api/json-api';
import { SchoolAccessService } from 'app/core/school-access/school-access.service';
import {
    Capability,
    CAPABILITY_LABELS,
    Department,
    GrantRequest,
    Membership,
    School,
} from 'app/core/school-access/school-access.types';
import { UserService } from 'app/core/user/user.service';
import { catchError, filter, forkJoin, map, of, switchMap } from 'rxjs';

const CAPABILITIES = Object.keys(CAPABILITY_LABELS) as Capability[];

@Component({
    selector: 'app-school-detail',
    standalone: true,
    imports: [
        DatePipe,
        FormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatSelectModule,
        MatSnackBarModule,
        MatTabsModule,
        RouterLink,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './school-detail.component.html',
})
export class SchoolDetailComponent {
    private readonly schoolAccess = inject(SchoolAccessService);
    private readonly adminUsers = inject(AdminUsersService);
    private readonly confirmation = inject(FuseConfirmationService);
    private readonly snackBar = inject(MatSnackBar);
    private readonly viewer = toSignal(inject(UserService).user$);
    readonly schoolID = inject(ActivatedRoute).snapshot.paramMap.get('schoolId') ?? '';

    readonly capabilities = CAPABILITIES;
    readonly labels = CAPABILITY_LABELS;
    readonly isSuperAdmin = computed(() => this.viewer()?.roles?.includes('SUPER_ADMIN') ?? false);

    readonly school = signal<School | null>(null);
    readonly departments = signal<Department[]>([]);
    /** Null while loading or when the viewer may not manage this school. */
    readonly memberships = signal<Membership[] | null>(null);
    readonly users = signal<AdminUser[]>([]);
    readonly canManage = signal(false);
    readonly loading = signal(true);
    readonly errorMessage = signal<string | null>(null);
    readonly busy = signal(false);

    readonly departmentName = signal('');
    readonly capabilityFilter = signal<Capability | ''>('');
    readonly departmentFilter = signal('');
    readonly grant = signal<{ userID: string; capability: Capability | ''; departmentID: string }>({
        userID: '',
        capability: '',
        departmentID: '',
    });

    readonly needsDepartment = computed(() => {
        const capability = this.grant().capability;
        return capability !== '' && capability !== 'manage_members';
    });
    readonly grantReady = computed(() => {
        const grant = this.grant();
        return !!grant.userID && !!grant.capability && (!this.needsDepartment() || !!grant.departmentID);
    });
    readonly visibleMemberships = computed(() =>
        (this.memberships() ?? []).filter(
            (m) =>
                (!this.capabilityFilter() || m.capability === this.capabilityFilter()) &&
                (!this.departmentFilter() || m.departmentID === this.departmentFilter())
        )
    );

    constructor() {
        this.load();
    }

    userName(id: string): string {
        const user = this.users().find((u) => u.id === id);
        return user ? `${user.name} (${user.email})` : id;
    }

    departmentLabel(id?: string): string {
        if (!id) return 'Whole school';
        return this.departments().find((d) => d.id === id)?.name ?? id;
    }

    updateGrant(change: Partial<{ userID: string; capability: Capability | ''; departmentID: string }>): void {
        this.grant.update((grant) => {
            const next = { ...grant, ...change };
            if (next.capability === 'manage_members') next.departmentID = '';
            return next;
        });
    }

    addDepartment(): void {
        const name = this.departmentName().trim();
        if (!name) return;
        this.busy.set(true);
        this.schoolAccess.createDepartment(this.schoolID, name).subscribe({
            next: (department) => {
                this.busy.set(false);
                this.departmentName.set('');
                this.departments.update((list) => [...list, department].sort((a, b) => a.name.localeCompare(b.name)));
                this.snackBar.open('Department added.', undefined, { duration: 3000 });
            },
            error: (error: HttpErrorResponse) => this.fail(error, 'Unable to add the department.'),
        });
    }

    submitGrant(): void {
        const grant = this.grant();
        if (!this.grantReady() || grant.capability === '') return;
        const request: GrantRequest = { userID: grant.userID, capability: grant.capability };
        if (this.needsDepartment()) request.departmentID = grant.departmentID;
        this.busy.set(true);
        this.schoolAccess.grant(this.schoolID, request).subscribe({
            next: () => {
                this.busy.set(false);
                this.grant.set({ userID: '', capability: '', departmentID: '' });
                this.snackBar.open('Capability granted.', undefined, { duration: 3000 });
                this.loadMemberships();
            },
            error: (error: HttpErrorResponse) => this.fail(error, 'Unable to grant the capability.'),
        });
    }

    revoke(membership: Membership): void {
        this.confirmation
            .open({
                title: 'Revoke capability?',
                message: `${this.userName(membership.userID)} loses “${this.labels[membership.capability]}” immediately, even if they are signed in. The record stays in the history.`,
                actions: { confirm: { label: 'Revoke' }, cancel: { label: 'Cancel' } },
            })
            .afterClosed()
            .pipe(
                filter((result) => result === 'confirmed'),
                switchMap(() => this.schoolAccess.revoke(this.schoolID, membership.id))
            )
            .subscribe({
                next: () => {
                    this.snackBar.open('Capability revoked.', undefined, { duration: 3000 });
                    this.loadMemberships();
                },
                error: (error: HttpErrorResponse) => this.fail(error, 'Unable to revoke the capability.'),
            });
    }

    private load(): void {
        forkJoin([
            this.schoolAccess.schools(),
            this.schoolAccess.departments(this.schoolID),
            this.schoolAccess.ownMemberships(),
        ]).subscribe({
            next: ([schools, departments, own]) => {
                this.school.set(schools.find((s) => s.id === this.schoolID) ?? null);
                this.departments.set(departments);
                this.canManage.set(
                    this.isSuperAdmin() ||
                        own.some((m) => m.schoolID === this.schoolID && m.capability === 'manage_members')
                );
                this.loading.set(false);
                if (this.canManage()) {
                    this.loadMemberships();
                    this.adminUsers
                        .query({ rows: 100, orderBy: 'name,ASC' })
                        .pipe(map((result) => result.items), catchError(() => of([])))
                        .subscribe((users) => this.users.set(users));
                }
            },
            error: (error: HttpErrorResponse) => {
                this.loading.set(false);
                this.fail(error, 'Unable to load this school.');
            },
        });
    }

    private loadMemberships(): void {
        this.schoolAccess.memberships(this.schoolID).subscribe({
            next: (memberships) => this.memberships.set(memberships),
            error: (error: HttpErrorResponse) => this.fail(error, 'Unable to load memberships.'),
        });
    }

    private fail(error: HttpErrorResponse, fallback: string): void {
        this.busy.set(false);
        this.errorMessage.set(
            error.status === 403
                ? 'You do not have permission to do that in this school.'
                : jsonApiErrorMessage(error as JsonApiErrorResponse, fallback)
        );
    }
}
