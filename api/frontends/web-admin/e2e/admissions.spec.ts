import { expect, Page, test } from '@playwright/test';

const staffUser = {
    id: '5cf37266-3473-4006-984f-9325122678b7',
    name: 'Admin Gopher',
    email: 'admin@example.com',
    roles: ['SCHOOL_ADMIN', 'ADMISSIONS_ADMIN', 'APPLICATION_REVIEWER'],
};

const applicantFixture = {
    applicationID: 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
    constituentID: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    programID: 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    academicTermID: 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
};

const adminEventFixture = {
    id: '9f5d8f95-5f4d-4d5e-8b97-e26dcb43bfb0',
    title: 'Kenya Admissions Open Day',
    type: 'open-day',
    status: 'upcoming',
    description:
        'Meet faculty, tour the campus, and review admissions options.',
    start: '2026-06-10T09:00:00Z',
    end: '2026-06-10T13:00:00Z',
    location: 'Main Campus Auditorium',
    isVirtual: false,
    capacity: 120,
    registrationDeadline: '2026-06-09T23:59:59Z',
    autoConfirmationEnabled: true,
    autoReminderEnabled: true,
};

const adminEventRegistrationsFixture = [
    {
        id: '1ad84d76-4d86-4ca9-aa4c-a7dfb8d973d1',
        constituentId: 'c_1',
        constituentName: 'Sofia Martinez',
        email: 'sofia.martinez@example.edu',
        phone: '+254700000001',
        status: 'registered',
        registeredAt: '2026-06-01T08:00:00Z',
        matchStatus: 'matched',
        source: 'portal',
    },
    {
        id: '2cb36f67-f7ee-4b55-b710-e4deed641469',
        constituentId: 'c_2',
        constituentName: 'James Okoro',
        email: 'james.okoro@example.edu',
        phone: '+254700000002',
        status: 'checked-in',
        registeredAt: '2026-06-01T09:00:00Z',
        matchStatus: 'matched',
        source: 'campaign',
        checkedInAt: '2026-06-10T08:55:00Z',
        checkedInById: staffUser.id,
    },
];

const applicantEventFixture = {
    id: '7dc3c4b1-6dd8-44c2-a255-d0c4e2f7de18',
    title: 'Applicant Open Day',
    type: 'open-day',
    status: 'upcoming',
    description: 'Join admissions advisors for an applicant-focused open day.',
    start: '2026-07-12T09:00:00Z',
    end: '2026-07-12T12:00:00Z',
    location: 'Nairobi Main Campus',
    isVirtual: false,
    capacity: 80,
    registeredCount: 12,
    checkedInCount: 0,
    registrationDeadline: '2026-07-11T21:00:00Z',
    autoConfirmationEnabled: true,
    autoReminderEnabled: true,
    registrations: [],
    dateCreated: '2026-06-01T00:00:00Z',
    dateUpdated: '2026-06-01T00:00:00Z',
};

test.describe('Admissions staff experience', () => {
    test.beforeEach(async ({ page }) => {
        await mockAuth(page);
        await signIn(page);
    });

    test('opens the applications list and application detail happy path', async ({
        page,
    }) => {
        await page.goto('/applications');

        await expect(
            page.getByRole('heading', { name: 'Applications' })
        ).toBeVisible();
        await expect(
            page.getByPlaceholder('Search applications...')
        ).toBeVisible();
        await expect(
            page.getByRole('button', { name: /New application/i })
        ).toBeVisible();

        await page.getByRole('link', { name: 'APP-3024' }).click();

        await expect(page).toHaveURL(/\/applications\/APP-3024$/);
        await expect(
            page.getByRole('link', { name: /Back to applications/i })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: /Achieng Otieno/i })
        ).toBeVisible();
        await expect(
            page.getByText('Application', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText('Documents', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'KUCCPS placement' })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'KCSE result snapshot' })
        ).toBeVisible();
    });

    test('shows the staff review workspace with authorized application queue', async ({
        page,
    }) => {
        await page.goto('/applications/review');

        await expect(
            page.getByText('Applications / Review Workspace')
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Staff review workspace' })
        ).toBeVisible();
        await expect(page.getByText('Authorized queue')).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Assigned applications' })
        ).toBeVisible();
        await expect(
            page.getByRole('button', { name: /Open application APP-3024/i })
        ).toBeVisible();
        await expect(page.getByText('Decision due')).toBeVisible();
    });

    test('navigates from lead pipeline to admissions settings', async ({
        page,
    }) => {
        await page.goto('/leads');

        await expect(page.getByText('Admissions / Leads')).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Lead pipeline' })
        ).toBeVisible();
        await expect(
            page.getByRole('combobox', { name: 'Score band' })
        ).toBeVisible();

        await page
            .getByRole('button', { name: 'Configure scoring rules' })
            .click();

        await expect(page).toHaveURL(/\/admissions-settings$/);
        await expect(page.getByText('Admissions / Settings')).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Settings & configuration' })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Custom field registry' })
        ).toBeVisible();
        await expect(
            page.getByText('Lead Scoring Rules', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText('Import/Export', { exact: true })
        ).toBeVisible();
    });

    test('dashboard renders KPIs and recent applications from API', async ({
        page,
    }) => {
        await page.goto('/dashboard');

        await expect(
            page.getByRole('heading', { name: /Welcome back/i })
        ).toBeVisible();
        await expect(page.getByText('Active applications')).toBeVisible();
        await expect(page.getByText('Submitted')).toBeVisible();
        await expect(page.getByText('Admitted')).toBeVisible();
        await expect(page.getByText('Enrolled')).toBeVisible();
        await expect(page.getByText('Recent applications')).toBeVisible();
        await expect(page.getByText('No applications yet.')).toBeVisible();
        await expect(page.getByText('Upcoming events')).toBeVisible();
    });

    test('loads admin events list, detail, and check-in flow from live APIs', async ({
        page,
    }) => {
        const eventRequests = await mockAdminEvents(page);

        await page.goto('/events');

        await expect(
            page.getByRole('heading', { name: 'Events' })
        ).toBeVisible();
        await expect(
            page.getByText('1 total events — live admissions schedule.')
        ).toBeVisible();
        await expect(
            page.getByRole('link', { name: adminEventFixture.title })
        ).toBeVisible();
        await expect(page.getByText('2 / 120')).toBeVisible();

        await page.getByRole('link', { name: adminEventFixture.title }).click();

        await expect(page).toHaveURL(
            new RegExp(`/events/${adminEventFixture.id}$`)
        );
        await expect(
            page.getByRole('heading', { name: adminEventFixture.title })
        ).toBeVisible();
        await page.getByRole('tab', { name: 'Registrations' }).click();
        await expect(page.getByText('Attendee List')).toBeVisible();
        await expect(page.getByText('Sofia Martinez')).toBeVisible();
        await page.getByRole('tab', { name: 'Check-in' }).click();
        await expect(
            page.getByRole('link', { name: 'Launch Check-in App' })
        ).toBeVisible();

        await page.getByRole('link', { name: 'Launch Check-in App' }).click();

        await expect(page).toHaveURL(
            new RegExp(`/events/${adminEventFixture.id}/checkin$`)
        );
        await expect(
            page.getByRole('heading', { name: adminEventFixture.title })
        ).toBeVisible();
        await expect(
            page.getByRole('banner').getByText('Checked In')
        ).toBeVisible();
        await expect(page.getByText('1 / 2')).toBeVisible();
        await page.getByRole('button', { name: 'Check In' }).click();
        await expect(page.getByText('2 / 2')).toBeVisible();
        await expect(eventRequests.checkIns).toBe(1);
    });
});

test.describe('Applicant portal experience', () => {
    test('starts the public application happy path', async ({ page }) => {
        const applicantRequests = await mockApplicantAdmissions(page);

        await page.goto('/portal');

        await expect(
            page.getByRole('link', { name: /SchoolCRM Kenya Applicant/i })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', {
                name: 'Start your admissions journey to School Of Thought University.',
            })
        ).toBeVisible();
        await expect(
            page.getByText(/KUCCPS placement or self-sponsored intake/i)
        ).toBeVisible();
        await expect(
            page.getByRole('link', { name: /Start application/i })
        ).toBeVisible();
        await expect(page.getByText('Email-only portal access')).toBeVisible();

        await page.getByRole('link', { name: /Start application/i }).click();

        await expect(page).toHaveURL(/\/portal\/apply$/);
        await expect(
            page.getByRole('heading', {
                name: 'Apply for the 2026 main intake',
            })
        ).toBeVisible();
        await expect(
            page.getByLabel('Application progress: step 1 of 6, account')
        ).toBeVisible();
        await expect(
            page.getByText('Applicant account', { exact: true })
        ).toBeVisible();
        await expect(page.getByLabel('Preferred programme')).not.toBeVisible();

        await page.getByLabel('First name').fill('John');
        await page.getByLabel('Last name').fill('Applicant');
        await page.getByLabel('Email').fill('applicant@example.com');
        await page.getByLabel('Phone').fill('+254712345678');
        await page.getByLabel('Date of birth').fill('2005-01-15');
        await page.getByLabel('Password', { exact: true }).fill('gophers');
        await page.getByLabel('Confirm password').fill('gophers');

        await expect(
            page.getByRole('button', { name: /Continue/i })
        ).toBeVisible();
        await expect(page.getByText('Application fee')).toBeVisible();
        await expect(page.getByText(/Ksh\s*150\.00/i)).toBeVisible();
        await expect(page.getByText(/Payment channel: M-Pesa/i)).toBeVisible();
        await expect(page.getByText('Need help?')).toBeVisible();
        await expect(
            page.getByText(/admissions@schoolcrm.ac.ke/i)
        ).toBeVisible();

        await page.getByRole('button', { name: /Continue/i }).click();
        await expect(
            page.getByText('Applicant account created.')
        ).toBeVisible();
        await expect(
            page.getByLabel('Application progress: step 2 of 6, programme')
        ).toBeVisible();
        await expect(
            page.getByText('Programme selection', { exact: true })
        ).toBeVisible();
        await expect(page.getByLabel('Preferred programme')).toBeVisible();
        await expect(page.getByLabel('Sponsorship route')).toBeVisible();
        await expect(page.getByLabel('Preferred programme')).toHaveValue(
            'Bachelor of Commerce'
        );
        await expect(page.getByLabel('Intake')).toHaveValue('2026 Main Intake');

        await page.getByRole('button', { name: /Continue/i }).click();
        await expect(
            page.getByText('Draft saved to admissions.')
        ).toBeVisible();
        await expect(
            page.getByLabel('Application progress: step 3 of 6, kcse details')
        ).toBeVisible();
        await expect(page.getByLabel('KCSE index number')).toBeVisible();
        await expect(page.getByLabel('KCSE mean grade')).toBeVisible();
        await page.getByLabel('KCSE index number').fill('12345678901/2025');
        await page.getByLabel('KCSE year').fill('2025');
        await page.getByLabel('KCSE mean grade').fill('B+');
        await page
            .getByLabel('Subject highlights')
            .fill('Maths B+, English A-');

        await page.getByRole('button', { name: /Documents/i }).click();
        await expect(
            page.getByLabel('Application progress: step 5 of 6, documents')
        ).toBeVisible();
        await expect(
            page.getByText('Supporting documents', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText('KCSE result document', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText('M-Pesa receipt', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText(/Drop KCSE result document or/i)
        ).toBeVisible();
        await expect(page.getByText(/Drop M-Pesa receipt or/i)).toBeVisible();

        await page.getByRole('button', { name: /Back/i }).click();
        await expect(
            page.getByLabel('Application progress: step 4 of 6, placement')
        ).toBeVisible();
        await expect(
            page.getByText('KUCCPS or self-sponsored placement', {
                exact: true,
            })
        ).toBeVisible();
        await expect(page.getByLabel('KUCCPS placement number')).toBeVisible();
        await expect(page.getByLabel('National ID or passport')).toBeVisible();
        await page
            .getByLabel('KUCCPS placement number')
            .fill('KUCCPS-2026-0001');
        await page.getByLabel('National ID or passport').fill('12345678');

        await page.getByRole('button', { name: /Review/i }).click();
        await expect(
            page.getByLabel('Application progress: step 6 of 6, review')
        ).toBeVisible();
        await expect(page.getByLabel('M-Pesa confirmation code')).toBeVisible();
        await expect(page.getByLabel('Review notes')).toBeVisible();
        await page.getByLabel('M-Pesa confirmation code').fill('QF123ABC45');
        await page.getByLabel('Review notes').fill('Ready for review');
        await page.getByRole('button', { name: /Submit application/i }).click();

        await expect(
            page.getByText('Application submitted for admissions review.')
        ).toBeVisible();
        expect(applicantRequests.onboardedApplicants).toBe(1);
        expect(
            applicantRequests.createdApplications +
                applicantRequests.savedApplications
        ).toBeGreaterThanOrEqual(1);
        expect(applicantRequests.submittedApplications).toBe(1);
    });

    test('shows Kenya-localized application status details', async ({
        page,
    }) => {
        await page.goto('/portal/status');

        await expect(
            page.getByRole('link', { name: /SchoolCRM Kenya Applicant/i })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Application status' })
        ).toBeVisible();
        await expect(page.getByText(/2026 Main Intake/i)).toBeVisible();
        await expect(
            page.getByText(/M-Pesa processing is not active in this preview/i)
        ).toBeVisible();
        await expect(page.getByText(/Ksh\s*150\.00/i)).toBeVisible();
        await expect(page.getByText('M-Pesa', { exact: true })).toBeVisible();
        await expect(
            page.getByText('KCSE result slip', { exact: true })
        ).toBeVisible();
        await expect(
            page.getByText(
                /Official KCSE result slip or certificate from KNEC/i
            )
        ).toBeVisible();
        await expect(page.getByText('KUCCPS placement letter')).toBeVisible();
        await expect(
            page.getByRole('heading', { name: 'Reach your admissions officer' })
        ).toBeVisible();
        await expect(
            page.getByText(/admissions@schoolcrm.ac.ke/i)
        ).toBeVisible();
    });

    test('loads applicant events and submits a portal registration', async ({
        page,
    }) => {
        const applicantRequests = await mockApplicantAdmissions(page);
        await seedApplicantPortalSession(page);

        await page.goto('/portal/events');

        await expect(
            page.getByRole('heading', { name: 'Visit, attend & apply' })
        ).toBeVisible();
        await expect(
            page.getByRole('heading', { name: applicantEventFixture.title })
        ).toBeVisible();
        await expect(page.getByText('68 spots left')).toBeVisible();

        await page.getByRole('button', { name: 'Reserve a spot' }).click();

        await expect(page.getByLabel('First name')).toHaveValue('John');
        await page.getByLabel('Last name').fill('Applicant');
        await page.getByLabel('Email').fill('applicant@example.com');
        await page.getByLabel('Phone').fill('+254712345678');
        await page.getByRole('button', { name: 'Submit registration' }).click();

        await expect(
            page.getByText('Registration confirmed for John Applicant.')
        ).toBeVisible();
        await expect(page.getByText('67 spots left')).toBeVisible();
        expect(applicantRequests.loadedEvents).toBeGreaterThanOrEqual(1);
        expect(applicantRequests.eventRegistrations).toBe(1);
    });

    test('submits a portal inquiry to admissions APIs', async ({ page }) => {
        const applicantRequests = await mockApplicantAdmissions(page);

        await page.goto(
            '/portal/inquiry?utm_source=google&utm_medium=cpc&utm_campaign=2026-intake'
        );

        await expect(
            page.getByRole('heading', { name: 'Ask admissions' })
        ).toBeVisible();
        await page.getByLabel('First name').fill('John');
        await page.getByLabel('Last name').fill('Applicant');
        await page.getByLabel('Date of birth').fill('2005-01-15');
        await page.getByLabel('Email').fill('applicant@example.com');
        await page.getByLabel('Phone').fill('+254712345678');
        await page
            .getByLabel('Programme of interest')
            .fill('Bachelor of Commerce');
        await page.getByLabel('Intake term').fill('2026 Main Intake');
        await page
            .getByLabel('Your question')
            .fill('What are the KCSE requirements for Commerce?');

        await page.getByRole('button', { name: 'Send inquiry' }).click();

        await expect(
            page.getByText('Your inquiry has been sent.')
        ).toBeVisible();
        expect(applicantRequests.submittedInquiries).toBe(1);
    });
});

async function mockApplicantAdmissions(page: Page): Promise<{
    onboardedApplicants: number;
    createdApplications: number;
    savedApplications: number;
    submittedApplications: number;
    loadedEvents: number;
    eventRegistrations: number;
    submittedInquiries: number;
}> {
    const requests = {
        onboardedApplicants: 0,
        createdApplications: 0,
        savedApplications: 0,
        submittedApplications: 0,
        loadedEvents: 0,
        eventRegistrations: 0,
        submittedInquiries: 0,
    };
    const applicantToken = createJwt({
        sub: 'f47ac10b-58cc-4372-a567-0e02b2c3d479',
        roles: ['APPLICANT'],
        portal: {
            scope: 'applicant_portal',
            applicationID: '',
            constituentID: applicantFixture.constituentID,
            email: 'applicant@example.com',
        },
    });

    await page.route('**/v1/auth/applicant-portal/onboard', async (route) => {
        requests.onboardedApplicants += 1;
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: {
                    id: applicantFixture.constituentID,
                    type: 'applicant-portal-token',
                    attributes: {
                        accessToken: applicantToken,
                        tokenType: 'Bearer',
                        expiresAt: new Date(
                            Date.now() + 3_600_000
                        ).toISOString(),
                        expiresIn: 3600,
                        applicationID: '',
                        constituentID: applicantFixture.constituentID,
                        applicantName: 'John Applicant',
                        email: 'applicant@example.com',
                    },
                },
            }),
        });
    });

    await page.route('**/v1/admissions/applicant/programs*', async (route) => {
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: [
                    {
                        id: applicantFixture.programID,
                        type: 'programs',
                        attributes: {
                            externalSISID: 'SIS-BCOM',
                            name: 'Bachelor of Commerce',
                            code: 'BCOM',
                            degreeLevel: 'UNDERGRADUATE',
                            active: true,
                            dateCreated: new Date().toISOString(),
                            dateUpdated: new Date().toISOString(),
                        },
                    },
                ],
                meta: { total: 1, page: 1, rowsPerPage: 100 },
            }),
        });
    });

    await page.route(
        '**/v1/admissions/applicant/academic-terms*',
        async (route) => {
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: [
                        {
                            id: applicantFixture.academicTermID,
                            type: 'academic-terms',
                            attributes: {
                                externalSISID: 'SIS-2026-MAIN',
                                name: '2026 Main Intake',
                                code: '2026-MAIN',
                                startDate: '2026-09-01T00:00:00Z',
                                endDate: '2027-04-30T00:00:00Z',
                                active: true,
                                dateCreated: new Date().toISOString(),
                                dateUpdated: new Date().toISOString(),
                            },
                        },
                    ],
                    meta: { total: 1, page: 1, rowsPerPage: 100 },
                }),
            });
        }
    );

    await page.route(
        '**/v1/admissions/applicant/application-form-templates*',
        async (route) => {
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: [
                        {
                            id: 'template-2026-main',
                            type: 'application-form-templates',
                            attributes: {
                                programID: applicantFixture.programID,
                                academicTermID: applicantFixture.academicTermID,
                                applicationType: 'KUCCPS_PLACEMENT',
                                name: '2026 Main Intake Form',
                                version: 1,
                                requiredFields: [],
                                checklistItems: [],
                                active: true,
                                priority: 1,
                                dateCreated: new Date().toISOString(),
                                dateUpdated: new Date().toISOString(),
                            },
                        },
                    ],
                    meta: { total: 1, page: 1, rowsPerPage: 100 },
                }),
            });
        }
    );

    await page.route('**/v1/admissions/applicant/events*', async (route) => {
        const request = route.request();
        const path = new URL(request.url()).pathname;

        if (request.method() !== 'GET' || path.includes('/registrations')) {
            await route.fallback();
            return;
        }

        requests.loadedEvents += 1;
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: [
                    {
                        id: applicantEventFixture.id,
                        type: 'events',
                        attributes: applicantEventFixture,
                    },
                ],
                meta: { total: 1, page: 1, rowsPerPage: 50 },
            }),
        });
    });

    await page.route(
        '**/v1/admissions/applicant/events/*/registrations',
        async (route) => {
            requests.eventRegistrations += 1;
            const payload = route.request().postDataJSON() as {
                firstName: string;
                lastName: string;
                email: string;
                phone?: string;
                source: string;
                matchStatus: string;
            };

            expect(payload).toMatchObject({
                firstName: 'John',
                lastName: 'Applicant',
                email: 'applicant@example.com',
                phone: '+254712345678',
                source: 'portal',
                matchStatus: 'matched',
            });

            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: {
                        id: '2c5da65c-b653-4f87-ac8c-6685d6bd6401',
                        type: 'event-registrations',
                        attributes: {
                            constituentId: applicantFixture.constituentID,
                            constituentName: `${payload.firstName} ${payload.lastName}`,
                            email: payload.email,
                            phone: payload.phone,
                            status: 'registered',
                            registeredAt: '2026-06-05T12:00:00Z',
                            matchStatus: payload.matchStatus,
                            source: payload.source,
                        },
                    },
                }),
            });
        }
    );

    await page.route('**/v1/admissions/inquiries', async (route) => {
        if (route.request().method() !== 'POST') {
            await route.fallback();
            return;
        }

        requests.submittedInquiries += 1;
        const payload = route.request().postDataJSON() as {
            firstName: string;
            lastName: string;
            dateOfBirth: string;
            primaryEmail: string;
            primaryPhone: string;
            programOfInterest?: string;
            termOfInterest?: string;
            source: string;
            utmSource?: string;
            utmMedium?: string;
            utmCampaign?: string;
            message?: string;
        };

        expect(payload).toMatchObject({
            firstName: 'John',
            lastName: 'Applicant',
            primaryEmail: 'applicant@example.com',
            primaryPhone: '+254712345678',
            programOfInterest: 'Bachelor of Commerce',
            termOfInterest: '2026 Main Intake',
            source: 'PORTAL_INQUIRY_FORM',
            utmSource: 'google',
            utmMedium: 'cpc',
            utmCampaign: '2026-intake',
            message: 'What are the KCSE requirements for Commerce?',
        });
        expect(payload.dateOfBirth).toMatch(/^2005-01-15T/);

        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: {
                    id: '4f5da65c-b653-4f87-ac8c-6685d6bd6402',
                    type: 'inquiries',
                    attributes: {
                        constituentID: applicantFixture.constituentID,
                        firstName: payload.firstName,
                        lastName: payload.lastName,
                        dateOfBirth: payload.dateOfBirth,
                        primaryEmail: payload.primaryEmail,
                        primaryPhone: payload.primaryPhone,
                        programOfInterest: payload.programOfInterest,
                        termOfInterest: payload.termOfInterest,
                        source: payload.source,
                        utmSource: payload.utmSource,
                        utmMedium: payload.utmMedium,
                        utmCampaign: payload.utmCampaign,
                        message: payload.message,
                        status: 'received',
                        dateCreated: new Date().toISOString(),
                        dateUpdated: new Date().toISOString(),
                    },
                },
            }),
        });
    });

    await page.route(
        '**/v1/admissions/applicant/applications',
        async (route) => {
            if (route.request().method() === 'POST') {
                requests.createdApplications += 1;
                await route.fulfill({
                    contentType: 'application/vnd.api+json',
                    body: applicantApplicationBody('DRAFT'),
                });
                return;
            }

            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: [],
                    meta: { total: 0, page: 1, rowsPerPage: 100 },
                }),
            });
        }
    );

    await page.route(
        '**/v1/admissions/applicant/applications/*/transitions',
        async (route) => {
            requests.submittedApplications += 1;
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: applicantApplicationBody('SUBMITTED'),
            });
        }
    );

    await page.route(
        '**/v1/admissions/applicant/applications/*',
        async (route) => {
            if (route.request().method() === 'PUT') {
                requests.savedApplications += 1;
            }

            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: applicantApplicationBody('DRAFT'),
            });
        }
    );

    return requests;
}

async function seedApplicantPortalSession(page: Page): Promise<void> {
    const applicantToken = createJwt({
        sub: 'f47ac10b-58cc-4372-a567-0e02b2c3d479',
        roles: ['APPLICANT'],
        portal: {
            scope: 'applicant_portal',
            applicationID: applicantFixture.applicationID,
            constituentID: applicantFixture.constituentID,
            email: 'applicant@example.com',
        },
    });

    await page.addInitScript(
        (session) => {
            localStorage.setItem(
                'applicantPortalSession',
                JSON.stringify(session)
            );
        },
        {
            id: applicantFixture.constituentID,
            accessToken: applicantToken,
            tokenType: 'Bearer',
            expiresAt: new Date(Date.now() + 3_600_000).toISOString(),
            expiresIn: 3600,
            applicationID: applicantFixture.applicationID,
            constituentID: applicantFixture.constituentID,
            applicantName: 'John Applicant',
            email: 'applicant@example.com',
        }
    );
}

async function mockAuth(page: Page): Promise<void> {
    const token = createJwt({ sub: staffUser.id, roles: staffUser.roles });

    await page.route('**/v1/admissions/**', async (route) => {
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: [],
                meta: { total: 0, page: 1, rows: 50 },
            }),
        });
    });

    await page.route('**/v1/auth/login', async (route) => {
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: {
                    type: 'auth-token',
                    attributes: {
                        accessToken: token,
                        tokenType: 'Bearer',
                        expiresAt: new Date(
                            Date.now() + 3_600_000
                        ).toISOString(),
                        expiresIn: 3600,
                        user: staffUser,
                    },
                },
            }),
        });
    });

    await page.route('**/v1/auth/authenticate', async (route) => {
        await route.fulfill({
            contentType: 'application/vnd.api+json',
            body: JSON.stringify({
                jsonapi: { version: '1.1' },
                data: {
                    type: 'auth-session',
                    attributes: {
                        userID: staffUser.id,
                        claims: { roles: staffUser.roles },
                        user: staffUser,
                    },
                },
            }),
        });
    });
}

async function mockAdminEvents(page: Page): Promise<{ checkIns: number }> {
    const state = {
        checkIns: 0,
        registrations: adminEventRegistrationsFixture.map((registration) => ({
            ...registration,
        })),
    };

    await page.route(/\/v1\/admissions\/events(?:\?|\/|$)/, async (route) => {
        const request = route.request();
        const url = new URL(request.url());
        const path = url.pathname;

        if (
            request.method() === 'GET' &&
            path.endsWith(`/events/${adminEventFixture.id}`)
        ) {
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: {
                        id: adminEventFixture.id,
                        type: 'events',
                        attributes: adminEventAttributes(state.registrations),
                    },
                }),
            });
            return;
        }

        if (
            request.method() === 'GET' &&
            path.endsWith(`/events/${adminEventFixture.id}/registrations`)
        ) {
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: state.registrations.map((registration) => ({
                        id: registration.id,
                        type: 'event-registrations',
                        attributes: {
                            ...registration,
                        },
                    })),
                    meta: {
                        total: state.registrations.length,
                        page: 1,
                        rowsPerPage: state.registrations.length,
                    },
                }),
            });
            return;
        }

        if (request.method() === 'GET' && path.endsWith('/events')) {
            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: [
                        {
                            id: adminEventFixture.id,
                            type: 'events',
                            attributes: adminEventAttributes(
                                state.registrations
                            ),
                        },
                    ],
                    meta: { total: 1, page: 1, rowsPerPage: 50 },
                }),
            });
            return;
        }

        await route.fallback();
    });

    await page.route(
        '**/v1/admissions/event-registrations/*/check-in',
        async (route) => {
            state.checkIns += 1;
            const registrationID =
                route.request().url().split('/').at(-2) ?? '';
            const registration = state.registrations.find(
                (item) => item.id === registrationID
            );

            if (registration) {
                registration.status = 'checked-in';
                registration.checkedInAt = '2026-06-10T09:15:00Z';
                registration.checkedInById = staffUser.id;
            }

            await route.fulfill({
                contentType: 'application/vnd.api+json',
                body: JSON.stringify({
                    jsonapi: { version: '1.1' },
                    data: {
                        id: registrationID,
                        type: 'event-registrations',
                        attributes: {
                            ...registration,
                        },
                    },
                }),
            });
        }
    );

    return state;
}

async function signIn(page: Page): Promise<void> {
    await page.goto('/sign-in');
    await page.getByLabel('Email address').fill(staffUser.email);
    await page.getByLabel('Password').fill('gophers');
    await page.getByRole('button', { name: 'Sign in' }).click();
    await expect(page).toHaveURL(/\/dashboard$/);
}

function applicantApplicationBody(status: 'DRAFT' | 'SUBMITTED'): string {
    return JSON.stringify({
        jsonapi: { version: '1.1' },
        data: {
            id: applicantFixture.applicationID,
            type: 'applications',
            attributes: {
                constituentID: applicantFixture.constituentID,
                programID: applicantFixture.programID,
                academicTermID: applicantFixture.academicTermID,
                applicationType: 'KUCCPS_PLACEMENT',
                status,
                kuccpsPlacement: {
                    placementID: 'KUCCPS-2026-0001',
                    institutionCode: 'STU',
                    programmeCode: 'BCOM',
                    programmeName: 'Bachelor of Commerce',
                    placementYear: 2025,
                },
                kcseResult: {
                    indexNumber: '12345678901/2025',
                    examYear: 2025,
                    subjects: [
                        { subjectCode: 'MATHS', grade: 'B+', points: 10 },
                        { subjectCode: 'ENGLISH', grade: 'A-', points: 11 },
                    ],
                    meanGrade: 'B+',
                    meanPoints: 10,
                },
                submittedAt:
                    status === 'SUBMITTED'
                        ? new Date().toISOString()
                        : undefined,
                dateCreated: new Date().toISOString(),
                dateUpdated: new Date().toISOString(),
            },
        },
    });
}

function adminEventAttributes(
    registrations: Array<{
        id: string;
        constituentId?: string;
        constituentName: string;
        email: string;
        phone?: string;
        status: string;
        registeredAt: string;
        matchStatus: string;
        source: string;
        checkedInAt?: string;
        checkedInById?: string;
    }>
): Record<string, unknown> {
    return {
        ...adminEventFixture,
        registeredCount: registrations.filter(
            (registration) => registration.status !== 'cancelled'
        ).length,
        checkedInCount: registrations.filter(
            (registration) => registration.status === 'checked-in'
        ).length,
        registrations,
        dateCreated: '2026-05-20T10:00:00Z',
        dateUpdated: '2026-06-05T10:00:00Z',
    };
}

function createJwt(payload: Record<string, unknown>): string {
    const now = Math.floor(Date.now() / 1000);
    const header = { alg: 'none', typ: 'JWT' };
    const claims = {
        iss: 'playwright',
        iat: now,
        exp: now + 3600,
        ...payload,
    };

    return `${base64Url(header)}.${base64Url(claims)}.`;
}

function base64Url(value: Record<string, unknown>): string {
    return Buffer.from(JSON.stringify(value))
        .toString('base64')
        .replace(/=/g, '')
        .replace(/\+/g, '-')
        .replace(/\//g, '_');
}
