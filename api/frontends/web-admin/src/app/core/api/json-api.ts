export interface JsonApiResource<T> {
    id: string;
    type: string;
    attributes: Omit<T, 'id'>;
}

export interface JsonApiDocument<T> {
    data: JsonApiResource<T>;
}

export interface JsonApiCollectionDocument<T> {
    data: JsonApiResource<T>[];
    meta?: JsonApiCollectionMeta;
}

export interface JsonApiCollectionMeta {
    total?: number;
    page?: number;
    rowsPerPage?: number;
}

export interface PaginatedResult<T> {
    items: T[];
    total: number;
    page: number;
    rowsPerPage: number;
}

export interface JsonApiError {
    detail?: string;
    title?: string;
}

export interface JsonApiErrorResponse {
    error?: {
        errors?: JsonApiError[];
        message?: string;
    };
}

export function unwrapJsonApiResource<T extends { id: string }>(
    document: JsonApiDocument<T>
): T {
    return {
        ...document.data.attributes,
        id: document.data.id,
    } as T;
}

export function unwrapJsonApiCollection<T extends { id: string }>(
    document: JsonApiCollectionDocument<T>
): PaginatedResult<T> {
    return {
        items: document.data.map((item) => ({
            ...item.attributes,
            id: item.id,
        }) as T),
        total: document.meta?.total ?? document.data.length,
        page: document.meta?.page ?? 1,
        rowsPerPage: document.meta?.rowsPerPage ?? document.data.length,
    };
}

export function jsonApiErrorMessage(
    error: JsonApiErrorResponse,
    fallback: string
): string {
    return (
        error.error?.errors?.[0]?.detail ??
        error.error?.errors?.[0]?.title ??
        error.error?.message ??
        fallback
    );
}
