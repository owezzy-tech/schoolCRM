import { Component, computed, inject } from '@angular/core';
import { NgClass } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { BreadcrumbComponent } from '@shared/components/breadcrumb/breadcrumb.component';

type SmartPageKind = 'dashboard' | 'list' | 'form' | 'profile' | 'calendar' | 'report' | 'settings' | 'workspace';

interface SmartPageData {
  title: string;
  section: string;
  kind?: SmartPageKind;
  icon?: string;
}

interface SmartPageAction {
  label: string;
  icon: string;
  tone: string;
}

interface SmartPageMetric {
  label: string;
  icon: string;
  value: string;
  tone: string;
}

@Component({
  selector: 'app-smart-demo-page',
  imports: [BreadcrumbComponent, NgClass, RouterLink],
  templateUrl: './smart-demo-page.component.html',
  styleUrl: './smart-demo-page.component.scss',
})
export class SmartDemoPageComponent {
  private readonly route = inject(ActivatedRoute);

  readonly page = computed<SmartPageData>(() => {
    const data = this.route.snapshot.data;
    const fallbackSegment = this.route.snapshot.pathFromRoot
      .flatMap((snapshot) => snapshot.url)
      .map((segment) => segment.path)
      .filter(Boolean)
      .at(-1);
    return {
      title: String(data['title'] ?? this.toTitle(fallbackSegment ?? 'Smart Page')),
      section: String(data['section'] ?? 'Smart School'),
      kind: this.isSmartPageKind(data['kind']) ? data['kind'] : 'workspace',
      icon: typeof data['icon'] === 'string' ? data['icon'] : 'dashboard_customize',
    };
  });

  readonly breadscrums = computed(() => [
    {
      title: this.page().title,
      items: [this.page().section],
      active: this.page().title,
    },
  ]);

  readonly actions = computed<SmartPageAction[]>(() => {
    const kind = this.page().kind ?? 'workspace';
    const actionMap: Record<SmartPageKind, SmartPageAction[]> = {
      dashboard: [
        { label: 'View analytics', icon: 'query_stats', tone: 'bg-blue' },
        { label: 'Export summary', icon: 'download', tone: 'bg-green' },
        { label: 'Configure widgets', icon: 'tune', tone: 'bg-purple' },
      ],
      list: [
        { label: 'Add record', icon: 'add_circle', tone: 'bg-blue' },
        { label: 'Filter list', icon: 'filter_alt', tone: 'bg-orange' },
        { label: 'Export table', icon: 'table_view', tone: 'bg-green' },
      ],
      form: [
        { label: 'Save draft', icon: 'save', tone: 'bg-blue' },
        { label: 'Validate fields', icon: 'rule', tone: 'bg-cyan' },
        { label: 'Submit workflow', icon: 'send', tone: 'bg-green' },
      ],
      profile: [
        { label: 'Open profile', icon: 'account_circle', tone: 'bg-blue' },
        { label: 'Review timeline', icon: 'timeline', tone: 'bg-purple' },
        { label: 'Update contact', icon: 'contact_mail', tone: 'bg-green' },
      ],
      calendar: [
        { label: 'Create event', icon: 'event_available', tone: 'bg-blue' },
        { label: 'Sync schedule', icon: 'sync', tone: 'bg-cyan' },
        { label: 'Print timetable', icon: 'print', tone: 'bg-orange' },
      ],
      report: [
        { label: 'Generate report', icon: 'summarize', tone: 'bg-blue' },
        { label: 'Download PDF', icon: 'picture_as_pdf', tone: 'bg-red' },
        { label: 'Share insights', icon: 'ios_share', tone: 'bg-green' },
      ],
      settings: [
        { label: 'Review access', icon: 'admin_panel_settings', tone: 'bg-blue' },
        { label: 'Update policy', icon: 'settings_suggest', tone: 'bg-purple' },
        { label: 'Audit logs', icon: 'fact_check', tone: 'bg-orange' },
      ],
      workspace: [
        { label: 'Open workspace', icon: 'open_in_new', tone: 'bg-blue' },
        { label: 'Attach documents', icon: 'attach_file', tone: 'bg-cyan' },
        { label: 'Notify users', icon: 'campaign', tone: 'bg-green' },
      ],
    };

    return actionMap[kind];
  });

  readonly metrics = computed<SmartPageMetric[]>(() => {
    const title = this.page().title;
    return [
      { label: `${title} queue`, icon: 'pending_actions', value: 'Ready', tone: 'l-bg-blue' },
      { label: 'Template coverage', icon: 'verified', value: 'Mapped', tone: 'l-bg-green' },
      { label: 'Workflow status', icon: 'bolt', value: 'Active', tone: 'l-bg-orange' },
      { label: 'Role access', icon: 'shield', value: this.page().section, tone: 'l-bg-purple' },
    ];
  });

  readonly detailCards = computed(() => [
    {
      heading: 'Smart template parity',
      body: `${this.page().title} is mapped from the SmartAngular v21 demo navigation and uses this app's existing card, breadcrumb, Material icon, and Bootstrap utility styles.`,
      icon: 'auto_awesome_mosaic',
    },
    {
      heading: 'School CRM workflow',
      body: 'Use this screen as the routed workspace for domain data, forms, tables, approvals, and reports while preserving the demo layout contract.',
      icon: 'school',
    },
    {
      heading: 'Responsive shell',
      body: 'The page follows the current fixed header/sidebar shell and collapses through the existing Smart SCSS media queries.',
      icon: 'devices',
    },
  ]);

  private toTitle(value: string): string {
    return value
      .split('-')
      .filter(Boolean)
      .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
      .join(' ');
  }

  private isSmartPageKind(value: unknown): value is SmartPageKind {
    return (
      value === 'dashboard' ||
      value === 'list' ||
      value === 'form' ||
      value === 'profile' ||
      value === 'calendar' ||
      value === 'report' ||
      value === 'settings' ||
      value === 'workspace'
    );
  }
}
