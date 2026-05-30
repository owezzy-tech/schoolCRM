import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';

interface AuditEvent {
    id: string;
    timestamp: string;
    actor: string;
    action: string;
    entity: string;
    ip: string;
}

@Component({
    selector: 'app-audit',
    standalone: true,
    imports: [MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './audit.component.html',
})
export class AuditComponent {
    readonly filters = ['Date range', 'Actor', 'Action'];

    readonly events: readonly AuditEvent[] = [
        { id: '1', timestamp: 'May 30, 10:45 AM', actor: 'Avery', action: 'Update', entity: 'Application (APP-3024)', ip: '192.168.1.105' },
        { id: '2', timestamp: 'May 30, 10:30 AM', actor: 'System', action: 'Create', entity: 'Communication (MSG-8492)', ip: '10.0.0.1' },
        { id: '3', timestamp: 'May 30, 09:15 AM', actor: 'Diego', action: 'Login', entity: 'Session', ip: '172.16.254.1' },
        { id: '4', timestamp: 'May 29, 16:20 PM', actor: 'Priya', action: 'Delete', entity: 'Duplicate Constituent (CON-921)', ip: '192.168.1.112' },
        { id: '5', timestamp: 'May 29, 14:05 PM', actor: 'System', action: 'Create', entity: 'Application (APP-3025)', ip: '10.0.0.1' },
        { id: '6', timestamp: 'May 29, 11:30 AM', actor: 'Avery', action: 'Update', entity: 'Campaign (CMP-42)', ip: '192.168.1.105' },
        { id: '7', timestamp: 'May 28, 09:00 AM', actor: 'Admin', action: 'Update', entity: 'Settings (Integration)', ip: '192.168.1.5' },
        { id: '8', timestamp: 'May 28, 08:45 AM', actor: 'Admin', action: 'Login', entity: 'Session', ip: '192.168.1.5' },
    ];
}
