export interface AuditEvent {
    id: string;
    obj_id: string;
    obj_domain: string;
    obj_name: string;
    actor_id: string;
    action: string;
    data: string;
    message: string;
    timestamp: string;
}

export interface AuditEventQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    obj_id?: string;
    obj_domain?: string;
    obj_name?: string;
    actor_id?: string;
    action?: string;
    since?: string;
    until?: string;
}
