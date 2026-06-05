import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NgxUiLoaderModule } from 'ngx-ui-loader';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss'],
    imports: [NgxUiLoaderModule, RouterOutlet],
})
export class AppComponent {
    /**
     * Constructor
     */
    constructor() {}
}
