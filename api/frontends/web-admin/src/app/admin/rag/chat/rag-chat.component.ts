import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { BreadcrumbComponent } from '@shared/components/breadcrumb/breadcrumb.component';
import { RagService } from '@core/service/rag.service';

type ChatAuthor = 'assistant' | 'user' | 'system';

interface ChatMessage {
  id: number;
  author: ChatAuthor;
  text: string;
  snippets: string[];
  documentIds: string[];
}

@Component({
  selector: 'app-rag-chat',
  templateUrl: './rag-chat.component.html',
  styleUrls: ['./rag-chat.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    BreadcrumbComponent,
    MatButtonModule,
    MatCardModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
    ReactiveFormsModule,
  ],
})
export class RagChatComponent {
  private readonly ragService = inject(RagService);
  private nextMessageId = 1;

  readonly breadscrums = [
    {
      title: 'RAG Chat',
      items: ['Admin'],
      active: 'RAG Chat',
    },
  ];

  readonly questionControl = new FormControl('', {
    nonNullable: true,
    validators: [Validators.required],
  });

  readonly uploadForm = new FormGroup({
    title: new FormControl('', {
      nonNullable: true,
      validators: [Validators.required],
    }),
    source: new FormControl('web-admin', {
      nonNullable: true,
      validators: [Validators.required],
    }),
  });

  readonly messages = signal<ChatMessage[]>([
    this.createMessage(
      'assistant',
      'Ask a question about indexed school documents, or attach a document to add it to the knowledge base.'
    ),
  ]);
  readonly selectedFile = signal<File | null>(null);
  readonly errorMessage = signal('');
  readonly uploadMessage = signal('');
  readonly isSending = signal(false);
  readonly isUploading = signal(false);
  readonly canSend = computed(() => this.questionControl.valid && !this.isSending());
  readonly canUpload = computed(
    () => this.uploadForm.valid && this.selectedFile() !== null && !this.isUploading()
  );

  onQuestionKeydown(event: KeyboardEvent): void {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      this.sendQuestion();
    }
  }

  sendQuestion(): void {
    const question = this.questionControl.value.trim();
    if (!question || this.isSending()) {
      this.questionControl.markAsTouched();
      return;
    }

    this.errorMessage.set('');
    this.uploadMessage.set('');
    this.isSending.set(true);
    this.messages.update((messages) => [...messages, this.createMessage('user', question)]);
    this.questionControl.reset('');

    this.ragService.query({ question, top_k: 3 }).subscribe({
      next: (response) => {
        this.messages.update((messages) => [
          ...messages,
          this.createMessage('assistant', response.answer, response.snippets, response.document_ids),
        ]);
        this.isSending.set(false);
      },
      error: () => {
        this.errorMessage.set('Could not query the RAG service. Please try again.');
        this.messages.update((messages) => [
          ...messages,
          this.createMessage('system', 'The RAG service did not return an answer.'),
        ]);
        this.isSending.set(false);
      },
    });
  }

  onFileSelected(event: Event): void {
    const input = event.target;
    if (!(input instanceof HTMLInputElement)) {
      return;
    }

    const file = input.files?.item(0) ?? null;
    this.selectedFile.set(file);
    this.uploadMessage.set('');
    this.errorMessage.set('');

    if (file && !this.uploadForm.controls.title.value.trim()) {
      this.uploadForm.controls.title.setValue(this.titleFromFile(file));
    }
  }

  uploadDocument(): void {
    const file = this.selectedFile();
    const title = this.uploadForm.controls.title.value.trim();
    const source = this.uploadForm.controls.source.value.trim();

    if (!file || !title || !source || this.isUploading()) {
      this.uploadForm.markAllAsTouched();
      return;
    }

    this.errorMessage.set('');
    this.uploadMessage.set('');
    this.isUploading.set(true);

    this.ragService.uploadDocument(title, source, file).subscribe({
      next: (response) => {
        this.uploadMessage.set(
          `Indexed ${response.chunk_count} chunks for document ${response.document_id}.`
        );
        this.messages.update((messages) => [
          ...messages,
          this.createMessage('system', `Document "${title}" is ${response.status}.`),
        ]);
        this.selectedFile.set(null);
        this.uploadForm.reset({ title: '', source: 'web-admin' });
        this.isUploading.set(false);
      },
      error: () => {
        this.errorMessage.set('Could not upload the selected document. Please try again.');
        this.isUploading.set(false);
      },
    });
  }

  trackMessage(_index: number, message: ChatMessage): number {
    return message.id;
  }

  private createMessage(
    author: ChatAuthor,
    text: string,
    snippets: string[] = [],
    documentIds: string[] = []
  ): ChatMessage {
    const message = {
      id: this.nextMessageId,
      author,
      text,
      snippets,
      documentIds,
    };
    this.nextMessageId += 1;
    return message;
  }

  private titleFromFile(file: File): string {
    return file.name.replace(/\.[^.]+$/, '').replace(/[-_]+/g, ' ');
  }
}
