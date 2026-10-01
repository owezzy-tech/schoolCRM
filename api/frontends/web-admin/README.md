# Fuse - Admin template and Starter project for Angular

This project was generated with [Angular CLI](https://github.com/angular/angular-cli)

## Development server

Run `ng serve` for a dev server. Navigate to `http://localhost:4400/`. The application will automatically reload if you change any of the source files.

## Storybook

From this frontend directory:

```bash
npm run storybook        # http://localhost:6006
npm run build-storybook  # static output in dist/storybook
```

From the repository root, use `npm run storybook` and `npm run build:storybook`.

Storybook uses the application's development build target for Fuse, Angular Material,
Tailwind and font styles. Public assets provide the same SVG icon sets as the app.
The theme toolbar switches between light and dark schemes.

Stories live in `src/stories/*.stories.ts`. The initial catalogue covers Fuse alerts
and cards, including review, publication, permission and curriculum layout examples.
These use fixture content and do not implement lesson storage or publishing. They
run without the auth, Go or RAG services. The dismissible alert includes a `play`
interaction that checks dismissal and its emitted event when that story is opened.

Use standalone presentational components with typed args and deterministic fixture
data. Supply only the providers needed by the component; avoid importing `app.config.ts`,
which starts the real application services. Storybook dependencies are development-only,
and its output is separate from the deployed application build.

## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory.

## Running unit tests

Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests

Run `ng e2e` to execute the end-to-end tests via a platform of your choice. To use this command, you need to first add a package that implements end-to-end testing capabilities.

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI Overview and Command Reference](https://angular.io/cli) page.
