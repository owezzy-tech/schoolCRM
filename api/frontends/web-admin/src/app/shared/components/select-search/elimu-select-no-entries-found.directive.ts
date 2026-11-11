/**
 * Forked from ngx-mat-select-search 9.0.0 (https://github.com/bithost-gmbh/ngx-mat-select-search,
 * tag 9.0.0, commit 9fa3f73) and renamed for SchoolCRM. Copyright (c) 2018 Bithost GmbH.
 * MIT licence: see LICENSE in this directory.
 */

import { Directive } from '@angular/core';

/**
 * Directive for providing a custom no entries found element.
 * e.g.
 * <elimu-select-search [formControl]="bankFilterCtrl">
 *   <span elimuSelectNoEntriesFound>
 *     No entries found <button>Add</button>
 *   </span>
 * </elimu-select-search>
 */
@Directive({
    selector: '[elimuSelectNoEntriesFound]',
})
export class ElimuSelectNoEntriesFoundDirective {}
