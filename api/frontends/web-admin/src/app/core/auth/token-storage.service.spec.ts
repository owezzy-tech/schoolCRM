import { TestBed } from '@angular/core/testing';
import { beforeEach, describe, expect, it } from 'vitest';

import { TokenStorageService } from './token-storage.service';

describe('TokenStorageService', () => {
    let service: TokenStorageService;

    beforeEach(() => {
        localStorage.clear();

        TestBed.configureTestingModule({});
        service = TestBed.inject(TokenStorageService);
    });

    it('stores and retrieves the access token', () => {
        service.set('access-token');

        expect(service.get()).toBe('access-token');
    });

    it('clears the access token', () => {
        service.set('access-token');
        service.clear();

        expect(service.get()).toBeNull();
    });
});
