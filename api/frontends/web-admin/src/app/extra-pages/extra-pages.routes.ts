import { Route } from "@angular/router";
import { BlankComponent } from "./blank/blank.component";
import { SmartDemoPageComponent } from "@shared/components/smart-demo-page/smart-demo-page.component";
export const EXTRA_PAGES_ROUTE: Route[] = [
  {
    path: "blank",
    component: BlankComponent,
  },
  {
    path: "**",
    component: SmartDemoPageComponent,
    data: { section: "Extra Pages", kind: "profile", icon: "description" },
  },
];
