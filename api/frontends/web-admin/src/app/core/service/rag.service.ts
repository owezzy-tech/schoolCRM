import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { environment } from 'environments/environment';
import { Observable } from 'rxjs';

import {
  RagDocumentDeleteResponse,
  RagDocumentUploadResponse,
  RagQueryRequest,
  RagQueryResponse,
} from '../models/rag';

@Injectable({
  providedIn: 'root',
})
export class RagService {
  private readonly http = inject(HttpClient);
  private readonly ragApiUrl = environment.ragApiUrl;

  query(request: RagQueryRequest): Observable<RagQueryResponse> {
    return this.http.post<RagQueryResponse>(`${this.ragApiUrl}/query`, request);
  }

  uploadDocument(
    title: string,
    source: string,
    file: File
  ): Observable<RagDocumentUploadResponse> {
    const formData = new FormData();
    formData.append('title', title);
    formData.append('source', source);
    formData.append('file', file);

    return this.http.post<RagDocumentUploadResponse>(`${this.ragApiUrl}/documents`, formData);
  }

  deleteDocument(documentId: string): Observable<RagDocumentDeleteResponse> {
    return this.http.delete<RagDocumentDeleteResponse>(`${this.ragApiUrl}/documents/${documentId}`);
  }
}
