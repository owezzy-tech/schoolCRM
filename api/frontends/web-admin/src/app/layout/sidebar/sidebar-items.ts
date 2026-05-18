import { RouteInfo } from './sidebar.metadata';

type SmartRole = '' | 'Admin' | 'Teacher' | 'Student' | 'All';

const specialTitles: Record<string, string> = {
  hr: 'Human Resources',
  faqs: 'FAQs',
  ui: 'User Interface',
  echart: 'ECharts',
  chartjs: 'Chart.js',
  'ngx-charts': 'ngx-charts',
  'ngx-datatable': 'ngx-datatable',
  'font-awesome': 'Font Awesome',
  page404: '404 - Not Found',
  page500: '500 - Server Error',
  first1: 'First',
  first2: 'Second',
  first3: 'Third',
};

const toTitle = (value: string): string => {
  const lastSegment = value.split('/').filter(Boolean).at(-1) ?? value;
  return (
    specialTitles[lastSegment] ??
    lastSegment
      .split('-')
      .filter(Boolean)
      .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
      .join(' ')
  );
};

const groupHeader = (title: string, role: SmartRole[] = ['All']): RouteInfo => ({
  path: '',
  title,
  iconType: '',
  icon: '',
  class: '',
  groupTitle: true,
  badge: '',
  badgeClass: '',
  role,
  submenu: [],
});

const leaf = (path: string, title = toTitle(path), role: SmartRole[] = ['']): RouteInfo => ({
  path,
  title,
  iconType: '',
  icon: '',
  class: 'ml-menu',
  groupTitle: false,
  badge: '',
  badgeClass: '',
  role,
  submenu: [],
});

const direct = (path: string, title: string, icon: string, role: SmartRole[]): RouteInfo => ({
  path,
  title,
  iconType: 'material-icons-outlined',
  icon,
  class: '',
  groupTitle: false,
  badge: '',
  badgeClass: '',
  role,
  submenu: [],
});

const menu = (title: string, icon: string, role: SmartRole[], paths: string[]): RouteInfo => ({
  path: '',
  title,
  iconType: 'material-icons-outlined',
  icon,
  class: 'menu-toggle',
  groupTitle: false,
  badge: '',
  badgeClass: '',
  role,
  submenu: paths.map((path) => leaf(path)),
});

const adminMenu = (title: string, icon: string, paths: string[]): RouteInfo =>
  menu(title, icon, ['Admin'], paths);

const teacherMenu = (title: string, icon: string, paths: string[]): RouteInfo =>
  menu(title, icon, ['Teacher'], paths);

const studentMenu = (title: string, icon: string, paths: string[]): RouteInfo =>
  menu(title, icon, ['Student'], paths);

export const ROUTES: RouteInfo[] = [
  groupHeader('MAIN'),
  adminMenu('Dashboard', 'space_dashboard', [
    '/admin/dashboard/main',
    '/admin/dashboard/dashboard2',
    '/admin/dashboard/teacher-dashboard',
    '/admin/dashboard/student-dashboard',
    '/admin/dashboard/library-dashboard',
    '/admin/dashboard/transport-dashboard',
  ]),
  adminMenu('Front Office', 'display_settings', [
    '/admin/front-office/admission-inquiry',
    '/admin/front-office/visitors',
    '/admin/front-office/complaints',
  ]),
  adminMenu('Teachers', 'person', [
    '/admin/teachers/all-teachers',
    '/admin/teachers/add-teacher',
    '/admin/teachers/edit-teacher',
    '/admin/teachers/about-teacher',
    '/admin/teachers/teacher-timetable',
    '/admin/teachers/assign-class-teacher',
  ]),
  adminMenu('Students', 'people_alt', [
    '/admin/students/all-students',
    '/admin/students/add-student',
    '/admin/students/edit-student',
    '/admin/students/student-attendance',
    '/admin/students/student-promotion',
    '/admin/students/student-certificates',
    '/admin/students/student-discipline',
    '/admin/students/student-health-records',
    '/admin/students/about-student',
  ]),
  adminMenu('Admissions', 'assignment_ind', [
    '/admin/admissions/admission-enquiries',
    '/admin/admissions/online-applications',
    '/admin/admissions/entrance-exams',
    '/admin/admissions/merit-list',
    '/admin/admissions/seat-allocation',
  ]),
  adminMenu('Examination', 'event_note', [
    '/admin/examination/exam-types',
    '/admin/examination/exam-schedule',
    '/admin/examination/hall-allocation',
    '/admin/examination/marks-entry',
    '/admin/examination/result-generation',
    '/admin/examination/report-cards',
  ]),
  adminMenu('Courses', 'school', [
    '/admin/courses/all-courses',
    '/admin/courses/add-course',
    '/admin/courses/edit-course',
    '/admin/courses/about-course',
  ]),
  adminMenu('Library', 'local_library', [
    '/admin/library/all-assets',
    '/admin/library/add-asset',
    '/admin/library/edit-asset',
    '/admin/library/book-status',
    '/admin/library/issue-return',
    '/admin/library/library-reports',
  ]),
  adminMenu('Staff', 'face', [
    '/admin/staff/all-staff',
    '/admin/staff/add-staff',
    '/admin/staff/edit-staff',
    '/admin/staff/about-staff',
    '/admin/staff/staff-attendance',
  ]),
  adminMenu('Human Resources', 'manage_accounts', [
    '/admin/human-resources/leave-requests',
    '/admin/human-resources/leave-balance',
    '/admin/human-resources/leave-types',
    '/admin/human-resources/holidays',
    '/admin/human-resources/todays-attendance',
    '/admin/human-resources/attendance-detail',
    '/admin/human-resources/attendance-sheet',
    '/admin/human-resources/employee-salary',
    '/admin/human-resources/payslip',
  ]),
  adminMenu('Transport', 'directions_bus', [
    '/admin/transport/vehicles',
    '/admin/transport/routes',
    '/admin/transport/drivers',
    '/admin/transport/student-allocation',
    '/admin/transport/transport-fees',
  ]),
  adminMenu('Holidays', 'airline_seat_individual_suite', [
    '/admin/holidays/all-holidays',
    '/admin/holidays/add-holiday',
    '/admin/holidays/edit-holiday',
  ]),
  adminMenu('Communication', 'campaign', [
    '/admin/communication/notice-board',
    '/admin/communication/announcements',
    '/admin/communication/sms-email',
  ]),
  direct('/admin/rag/chat', 'RAG Chat', 'chat', ['Admin']),
  adminMenu('Fees', 'monetization_on', [
    '/admin/fees/all-fees',
    '/admin/fees/add-fees',
    '/admin/fees/edit-fees',
    '/admin/fees/fees-type',
    '/admin/fees/fees-discount',
    '/admin/fees/fee-receipt',
  ]),
  adminMenu('Class', 'view_comfy_alt', ['/admin/class/class-list', '/admin/class/class-timetable']),
  adminMenu('Academics', 'auto_stories', [
    '/admin/academics/academic-year',
    '/admin/academics/sessions',
    '/admin/academics/classes',
    '/admin/academics/subjects',
    '/admin/academics/course-curriculum',
    '/admin/academics/assignment',
    '/admin/academics/lesson-planning',
  ]),
  adminMenu('Hostel', 'gite', [
    '/admin/hostel/room-list',
    '/admin/hostel/room-type',
    '/admin/hostel/allocations',
    '/admin/hostel/attendance',
    '/admin/hostel/hostel-fees',
  ]),
  adminMenu('Departments', 'business', [
    '/admin/departments/all-departments',
    '/admin/departments/add-department',
    '/admin/departments/edit-department',
  ]),
  adminMenu('Reports', 'assessment', [
    '/admin/reports/academic-reports',
    '/admin/reports/attendance-reports',
    '/admin/reports/fee-reports',
    '/admin/reports/exam-reports',
    '/admin/reports/custom-reports',
  ]),
  adminMenu('Settings', 'settings', [
    '/admin/settings/institute-profile',
    '/admin/settings/role-permissions',
    '/admin/settings/user-management',
    '/admin/settings/academic-rules',
    '/admin/settings/notification-settings',
    '/admin/settings/system-logs',
    '/admin/settings/backup-restore',
  ]),

  groupHeader('TEACHER', ['Teacher']),
  direct('/teacher/dashboard', 'Dashboard', 'space_dashboard', ['Teacher']),
  direct('/teacher/today-schedule', 'Today Schedule', 'access_time', ['Teacher']),
  teacherMenu('Academics', 'auto_stories', [
    '/teacher/academics/my-classes',
    '/teacher/academics/my-subjects',
    '/teacher/academics/lesson-plans',
    '/teacher/academics/study-materials',
    '/teacher/academics/assignments',
  ]),
  teacherMenu('Students', 'people_alt', [
    '/teacher/students/class-students',
    '/teacher/students/student-profiles',
    '/teacher/students/student-attendance',
    '/teacher/students/student-performance',
  ]),
  teacherMenu('Examination', 'history_edu', [
    '/teacher/examination/exam-schedule',
    '/teacher/examination/marks-entry',
    '/teacher/examination/grade-submission',
    '/teacher/examination/result-preview',
  ]),
  teacherMenu('Timetable', 'table_chart', [
    '/teacher/timetable/my-timetable',
    '/teacher/timetable/substitution-requests',
  ]),
  teacherMenu('Communication', 'chat', [
    '/teacher/communication/notices',
    '/teacher/communication/messages',
    '/teacher/communication/announcements',
  ]),
  teacherMenu('Attendance', 'how_to_reg', [
    '/teacher/attendance/daily-attendance',
    '/teacher/attendance/attendance-summary',
  ]),
  teacherMenu('Leave', 'offline_pin', ['/teacher/leave/apply-leave', '/teacher/leave/leave-status']),
  direct('/teacher/lectures', 'Lectures', 'menu_book', ['Teacher']),
  teacherMenu('My Profile', 'assignment_ind', [
    '/teacher/my-profile/profile-info',
    '/teacher/my-profile/documents',
    '/teacher/my-profile/change-password',
  ]),
  direct('/teacher/settings', 'Settings', 'settings', ['Teacher']),

  groupHeader('STUDENT', ['Student']),
  direct('/student/dashboard', 'Dashboard', 'space_dashboard', ['Student']),
  direct('/student/homework', 'Homework', 'article', ['Student']),
  direct('/student/timetable', 'Timetable', 'table_chart', ['Student']),
  direct('/student/my-class', 'My Class', 'co_present', ['Student']),
  studentMenu('Academics', 'school', [
    '/student/academics/syllabus',
    '/student/academics/assignments',
    '/student/academics/study-materials',
    '/student/academics/academic-calendar',
    '/student/academics/my-subjects',
  ]),
  studentMenu('Examination', 'assignment', [
    '/student/examination/exam-schedule',
    '/student/examination/hall-ticket',
    '/student/examination/marks',
    '/student/examination/report-card',
    '/student/examination/results',
  ]),
  studentMenu('Attendance', 'event_available', [
    '/student/attendance/monthly-summary',
    '/student/attendance/my-attendance',
  ]),
  direct('/student/leave-request', 'Leave Request', 'offline_pin', ['Student']),
  studentMenu('Fees', 'monetization_on', [
    '/student/fees/fee-receipts',
    '/student/fees/due-fees',
    '/student/fees/fee-details',
    '/student/fees/online-payment',
  ]),
  studentMenu('Library', 'local_library', [
    '/student/library/my-issued-books',
    '/student/library/due-dates',
    '/student/library/book-history',
  ]),
  studentMenu('Transport', 'directions_bus', [
    '/student/transport/my-route',
    '/student/transport/vehicle-details',
  ]),
  studentMenu('Hostel', 'hotel', [
    '/student/hostel/room-details',
    '/student/hostel/hostel-fees',
    '/student/hostel/complaints',
  ]),
  direct('/student/notices', 'Notices', 'announcement', ['Student']),
  direct('/student/profile', 'Profile', 'assignment_ind', ['Student']),
  direct('/student/settings', 'Settings', 'settings', ['Student']),

  groupHeader('APPS & UI', ['Admin']),
  direct('/calendar', 'Calendar', 'event_note', ['Admin']),
  direct('/task', 'Tasks', 'fact_check', ['Admin']),
  direct('/contacts', 'Contacts', 'contacts', ['Admin']),
  adminMenu('Email', 'email', ['/email/inbox', '/email/compose', '/email/read-mail']),
  adminMenu('More Apps', 'stars', ['/apps/chat', '/apps/dragdrop', '/apps/contact-grid', '/apps/support']),
  adminMenu('Widgets', 'widgets', ['/widget/chart-widget', '/widget/data-widget']),
  adminMenu('User Interface', 'dvr', [
    '/ui/alerts',
    '/ui/badges',
    '/ui/chips',
    '/ui/modal',
    '/ui/buttons',
    '/ui/expansion-panel',
    '/ui/bottom-sheet',
    '/ui/dialogs',
    '/ui/cards',
    '/ui/labels',
    '/ui/list-group',
    '/ui/snackbar',
    '/ui/preloaders',
    '/ui/progressbars',
    '/ui/tabs',
    '/ui/typography',
    '/ui/helper-classes',
  ]),
  adminMenu('Forms', 'subtitles', [
    '/forms/form-controls',
    '/forms/advance-controls',
    '/forms/form-example',
    '/forms/form-validation',
    '/forms/wizard',
    '/forms/editors',
  ]),
  adminMenu('Tables', 'view_list', ['/tables/basic-tables', '/tables/material-tables', '/tables/ngx-datatable']),
  adminMenu('Charts', 'insert_chart', [
    '/charts/echart',
    '/charts/apex',
    '/charts/chartjs',
    '/charts/ngx-charts',
    '/charts/gauge',
  ]),
  adminMenu('Timeline', 'timeline', ['/timeline/timeline1', '/timeline/timeline2']),
  adminMenu('Icons', 'eco', ['/icons/material', '/icons/font-awesome']),
  adminMenu('Authentication', 'supervised_user_circle', [
    '/authentication/signin',
    '/authentication/signup',
    '/authentication/forgot-password',
    '/authentication/locked',
    '/authentication/page404',
    '/authentication/page500',
  ]),
  adminMenu('Extra Pages', 'description', [
    '/extra-pages/profile',
    '/extra-pages/pricing',
    '/extra-pages/invoice',
    '/extra-pages/faqs',
    '/extra-pages/blank',
  ]),
  {
    path: '',
    title: 'Multi level Menu',
    iconType: 'material-icons-outlined',
    icon: 'slideshow',
    class: 'menu-toggle',
    groupTitle: false,
    badge: '',
    badgeClass: '',
    role: ['Admin'],
    submenu: [
      leaf('/multilevel/first1', 'First'),
      {
        ...leaf('/', 'Second'),
        class: 'ml-sub-menu',
        submenu: [
          { ...leaf('/multilevel/secondlevel/second1', 'Second 1'), class: 'ml-menu2' },
          {
            ...leaf('/', 'Second 2'),
            class: 'ml-sub-menu2',
            submenu: [{ ...leaf('/multilevel/thirdlevel/third1', 'Third 1'), class: 'ml-menu3' }],
          },
        ],
      },
      leaf('/multilevel/first3', 'Third'),
    ],
  },
];
