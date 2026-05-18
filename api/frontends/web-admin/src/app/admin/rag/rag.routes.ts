import { Route } from '@angular/router';

import { RagChatComponent } from './chat/rag-chat.component';

export const RAG_ROUTE: Route[] = [
  {
    path: '',
    redirectTo: 'chat',
    pathMatch: 'full',
  },
  {
    path: 'chat',
    component: RagChatComponent,
  },
];
