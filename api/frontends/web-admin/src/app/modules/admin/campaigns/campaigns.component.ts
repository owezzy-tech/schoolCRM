import { TitleCasePipe } from '@angular/common';
import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { Campaign } from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { catchError, of } from 'rxjs';

@Component({
    selector: 'app-campaigns',
    standalone: true,
    imports: [MatButtonModule, MatIconModule, TitleCasePipe],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './campaigns.component.html',
})
export class CampaignsComponent {
    private readonly admissionsService = inject(AdmissionsService);

    readonly loading = signal(false);
    readonly errorMessage = signal('');
    readonly campaigns = signal<Campaign[]>([]);
    readonly selectedCampaignId = signal('');
    readonly selectedCampaign = computed(
        () =>
            this.campaigns().find(
                (campaign) => campaign.id === this.selectedCampaignId()
            ) ?? this.campaigns()[0]
    );
    readonly totalAudience = computed(() =>
        this.campaigns().reduce(
            (total, campaign) => total + campaign.metrics.audienceSize,
            0
        )
    );
    readonly activeCampaigns = computed(
        () =>
            this.campaigns().filter((campaign) => campaign.status === 'ACTIVE')
                .length
    );
    readonly averageOpenRate = computed(() => {
        const sentCampaigns = this.campaigns().filter(
            (campaign) => campaign.metrics.sent > 0
        );

        if (!sentCampaigns.length) {
            return 0;
        }

        const openRate = sentCampaigns.reduce(
            (total, campaign) => total + this.rate(campaign.metrics.opened, campaign.metrics.sent),
            0
        );

        return Math.round(openRate / sentCampaigns.length);
    });

    constructor() {
        this.loadCampaigns();
    }

    loadCampaigns(): void {
        this.loading.set(true);
        this.errorMessage.set('');

        this.admissionsService
            .queryCampaigns({ page: 1, rows: 50, orderBy: 'starts_at,DESC' })
            .pipe(
                catchError((error) => {
                    this.errorMessage.set(
                        jsonApiErrorMessage(error, 'Unable to load campaigns.')
                    );
                    this.campaigns.set([]);

                    return of(undefined);
                })
            )
            .subscribe((result) => {
                this.loading.set(false);

                if (!result) {
                    return;
                }

                this.campaigns.set(result.items);

                if (!this.selectedCampaignId() && result.items.length) {
                    this.selectedCampaignId.set(result.items[0].id);
                }
            });
    }

    selectCampaign(campaignId: string): void {
        this.selectedCampaignId.set(campaignId);
    }

    statusClass(campaign: Campaign): string {
        switch (campaign.status) {
            case 'ACTIVE':
                return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300';
            case 'PAUSED':
                return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300';
            case 'COMPLETED':
                return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300';
            case 'DRAFT':
                return 'bg-slate-100 text-secondary dark:bg-slate-800';
        }
    }

    channelClass(campaign: Campaign): string {
        return campaign.channel === 'EMAIL'
            ? 'bg-primary-50 text-primary dark:bg-primary-900/20'
            : 'bg-purple-50 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300';
    }

    rate(value: number, total: number): number {
        if (!total) {
            return 0;
        }

        return Math.round((value / total) * 100);
    }

    compactNumber(value: number): string {
        return new Intl.NumberFormat('en-US').format(value);
    }

    schedule(campaign: Campaign): string {
        if (!campaign.startsAt && !campaign.endsAt) {
            return 'No schedule set';
        }

        return [campaign.startsAt, campaign.endsAt]
            .filter((value): value is string => Boolean(value))
            .map((value) => new Intl.DateTimeFormat('en-US').format(new Date(value)))
            .join(' - ');
    }
}
