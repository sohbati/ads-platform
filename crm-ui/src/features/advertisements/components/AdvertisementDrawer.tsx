"use client";

import { useEffect, useState } from "react";
import type { DescriptionsProps } from "antd";
import {
  Alert,
  App,
  Button,
  Descriptions,
  Drawer,
  Empty,
  Image,
  Space,
  Typography,
} from "antd";
import { AppError } from "@/components/ui/AppError";
import { AppLoading } from "@/components/ui/AppLoading";
import { ConfirmDialog } from "@/components/shared/ConfirmDialog";
import { mediaUrl } from "@/config/publicEnv";
import { can } from "@/lib/auth/permissions";
import { useSession } from "@/lib/auth/SessionProvider";
import { formatDateTime } from "@/lib/formatters/date";
import { formatDigits } from "@/lib/formatters/digits";
import { formatMobile } from "@/lib/formatters/phone";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { Permission } from "@/lib/permissions/codes";
import {
  AdvertisementStatus,
  CLAIM_HEARTBEAT_MS,
  CLAIM_WARNING_MS,
} from "../constants/advertisement.constants";
import {
  useAdvertisement,
  useClaimActions,
  useReviewMutations,
} from "../hooks/useAdvertisements";
import type { AdvertisementDetail } from "../types/advertisement.types";
import { AdvertisementStatusTag } from "./AdvertisementStatusTag";
import { formatPrice } from "./formatPrice";
import { RejectDialog } from "./RejectDialog";

interface AdvertisementDrawerProps {
  adId: number | null;
  onClose: () => void;
  /** Called after the reviewer approves or rejects, e.g. to open the next ad. */
  onDecided?: () => void;
}

function displayValue(value: unknown): string {
  if (value === null || value === undefined || value === "") {
    return "—";
  }
  if (typeof value === "boolean") {
    return value ? "✓" : "✗";
  }
  if (typeof value === "object") {
    return JSON.stringify(value);
  }
  return String(value);
}

/** The current time, refreshed every `intervalMs` while enabled. */
function useNow(intervalMs: number, enabled: boolean) {
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    if (!enabled) {
      return;
    }
    const timer = setInterval(() => setNow(Date.now()), intervalMs);
    return () => clearInterval(timer);
  }, [intervalMs, enabled]);
  return now;
}

export function AdvertisementDrawer({ adId, onClose, onDecided }: AdvertisementDrawerProps) {
  const { t, locale } = useI18n();
  const { message } = App.useApp();
  const session = useSession();
  const query = useAdvertisement(adId);
  const { approve, reject } = useReviewMutations();
  const { claim, extend, release, forceRelease } = useClaimActions();
  const [approving, setApproving] = useState(false);
  const [rejecting, setRejecting] = useState(false);
  const [freeing, setFreeing] = useState(false);
  const ad = query.data;

  const pending = ad?.status === AdvertisementStatus.pending;
  const claimedByMe = pending && ad?.claimed_by_staff_id === session.staffId;
  const now = useNow(5_000, claimedByMe);
  const remaining = ad?.claim_expires_at ? Date.parse(ad.claim_expires_at) - now : 0;
  const mine = claimedByMe && remaining > 0;
  const heldByOther = pending && ad?.claimed_by_staff_id != null && !claimedByMe;

  const canReview = pending && can(session.permissions, Permission.advertisementReview);
  const canApprove = mine && can(session.permissions, Permission.advertisementApprove);
  const canReject = mine && can(session.permissions, Permission.advertisementReject);
  const canFree =
    heldByOther && can(session.permissions, Permission.advertisementReviewManage);

  // Keep the claim alive while the reviewer has the ad open.
  const adIdHeld = mine ? ad?.id : undefined;
  const { mutate: extendClaim } = extend;
  useEffect(() => {
    if (adIdHeld === undefined) {
      return;
    }
    const timer = setInterval(() => extendClaim(adIdHeld), CLAIM_HEARTBEAT_MS);
    return () => clearInterval(timer);
  }, [adIdHeld, extendClaim]);

  // The claim ran out (idle too long): say so once and show the ad as it is now.
  const expiredClaim = claimedByMe && remaining <= 0 ? ad?.claim_expires_at : null;
  const { refetch } = query;
  useEffect(() => {
    if (expiredClaim) {
      message.warning(t("claimExpired"));
      void refetch();
    }
  }, [expiredClaim, message, refetch, t]);

  const close = () => {
    // Closing the ad hands it back, so the next reviewer does not wait for the claim to expire.
    if (mine && ad) {
      release.mutate(ad.id);
    }
    onClose();
  };

  const decided = () => {
    setApproving(false);
    setRejecting(false);
    onDecided?.();
  };

  return (
    <Drawer
      open={adId !== null}
      title={ad ? ad.title : t("adDetails")}
      size="large"
      onClose={close}
      extra={
        <Space>
          {mine ? <Button onClick={close}>{t("returnToQueue")}</Button> : null}
          {canReject ? (
            <Button danger onClick={() => setRejecting(true)}>
              {t("reject")}
            </Button>
          ) : null}
          {canApprove ? (
            <Button type="primary" onClick={() => setApproving(true)}>
              {t("approve")}
            </Button>
          ) : null}
          {canReview && !mine && !heldByOther ? (
            <Button
              type="primary"
              loading={claim.isPending}
              onClick={() => ad && claim.mutate(ad.id)}
            >
              {t("reviewThisAd")}
            </Button>
          ) : null}
          {canFree ? (
            <Button danger onClick={() => setFreeing(true)}>
              {t("forceRelease")}
            </Button>
          ) : null}
        </Space>
      }
    >
      {query.isLoading ? <AppLoading /> : null}
      {query.error ? (
        <AppError error={query.error} onRetry={() => void query.refetch()} />
      ) : null}
      {mine && remaining <= CLAIM_WARNING_MS ? (
        <Alert
          type="warning"
          showIcon
          style={{ marginBottom: 16 }}
          title={t("claimExpiresSoon")}
          action={
            <Button
              size="small"
              loading={claim.isPending}
              onClick={() => ad && claim.mutate(ad.id)}
            >
              {t("keepReviewing")}
            </Button>
          }
        />
      ) : mine ? (
        <Alert type="info" showIcon style={{ marginBottom: 16 }} title={t("reviewingByYou")} />
      ) : heldByOther && ad ? (
        <Alert
          type="info"
          showIcon
          style={{ marginBottom: 16 }}
          title={`${t("reviewing")}: ${ad.claimed_by_name ?? "—"} · ${formatDateTime(ad.claimed_at, locale)}`}
        />
      ) : null}
      {ad ? <Details ad={ad} locale={locale} t={t} /> : null}

      <ConfirmDialog
        open={approving}
        title={t("confirmApprove")}
        loading={approve.isPending}
        onCancel={() => setApproving(false)}
        onConfirm={() => {
          if (ad) {
            approve.mutate(ad.id, {
              onSuccess: decided,
              onError: () => setApproving(false),
            });
          }
        }}
      />
      <RejectDialog
        open={rejecting}
        loading={reject.isPending}
        onCancel={() => setRejecting(false)}
        onSubmit={(input) => {
          if (ad) {
            reject.mutate(
              { id: ad.id, input },
              { onSuccess: decided, onError: () => setRejecting(false) },
            );
          }
        }}
      />
      <ConfirmDialog
        open={freeing}
        title={t("confirmForceRelease")}
        loading={forceRelease.isPending}
        onCancel={() => setFreeing(false)}
        onConfirm={() => {
          if (ad) {
            forceRelease.mutate(ad.id, { onSettled: () => setFreeing(false) });
          }
        }}
      />
    </Drawer>
  );
}

function Details({
  ad,
  locale,
  t,
}: {
  ad: AdvertisementDetail;
  locale: "fa" | "en";
  t: ReturnType<typeof useI18n>["t"];
}) {
  const media = (ad.media ?? []).filter((item) => item.url || item.thumb);
  const keyValueItems = (entries: [string, unknown][]): DescriptionsProps["items"] =>
    entries.map(([key, value]) => ({
      key,
      label: <span dir="ltr">{key}</span>,
      children: <span dir="auto">{displayValue(value)}</span>,
    }));
  const attrs = keyValueItems(Object.entries(ad.attrs ?? {}));
  const contact = keyValueItems(Object.entries(ad.contact ?? {}));

  const summary: DescriptionsProps["items"] = [
    { key: "status", label: t("status"), children: <AdvertisementStatusTag ad={ad} /> },
    { key: "id", label: t("adId"), children: formatDigits(ad.id, locale) },
    {
      key: "category",
      label: t("category"),
      children: ad.category_title || formatDigits(ad.category_id, locale),
    },
    {
      key: "city",
      label: t("city"),
      children: `${ad.city_name || formatDigits(ad.city_id, locale)}${
        ad.location?.neighborhood ? `، ${ad.location.neighborhood}` : ""
      }`,
    },
    { key: "price", label: t("price"), children: formatPrice(ad, locale, t) },
    {
      key: "customer",
      label: t("customer"),
      children: (
        <>
          {ad.customer_name}{" "}
          <Typography.Text type="secondary" dir="ltr">
            {formatMobile(ad.customer_mobile, locale)}
          </Typography.Text>
        </>
      ),
    },
    {
      key: "created",
      label: t("createdAt"),
      children: formatDateTime(ad.created_at, locale),
    },
    ...(ad.pending_since
      ? [
          {
            key: "pending",
            label: t("waitingSince"),
            children: formatDateTime(ad.pending_since, locale),
          },
        ]
      : []),
  ];

  const reasons: DescriptionsProps["items"] = [
    {
      key: "customer",
      label: t("customerReason"),
      children: <span className="app-pre-wrap">{ad.rejection_reason_customer}</span>,
    },
    {
      key: "internal",
      label: t("internalReason"),
      children: (
        <span className="app-pre-wrap">{ad.rejection_reason_internal ?? "—"}</span>
      ),
    },
  ];

  return (
    <Space orientation="vertical" size="large" style={{ width: "100%" }}>
      {ad.reviewed_at ? (
        <Alert
          type={ad.rejection_reason_customer ? "warning" : "success"}
          showIcon
          title={`${t("lastReview")}: ${ad.reviewed_by_name ?? "—"} · ${formatDateTime(ad.reviewed_at, locale)}`}
          description={
            ad.rejection_reason_customer ? (
              <Descriptions size="small" column={1} items={reasons} />
            ) : null
          }
        />
      ) : null}

      <Descriptions bordered size="small" column={1} items={summary} />

      <section>
        <Typography.Title level={5}>{t("description")}</Typography.Title>
        <Typography.Paragraph className="app-pre-wrap">
          {ad.description}
        </Typography.Paragraph>
      </section>

      <section>
        <Typography.Title level={5}>{t("photos")}</Typography.Title>
        {media.length === 0 ? (
          <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={t("noPhotos")} />
        ) : (
          <Image.PreviewGroup>
            <Space wrap>
              {media.map((item, index) => (
                <Image
                  key={item.url ?? item.thumb ?? index}
                  width={120}
                  height={90}
                  style={{ objectFit: "cover", borderRadius: 6 }}
                  src={mediaUrl(item.thumb ?? item.url)}
                  preview={{ src: mediaUrl(item.url ?? item.thumb) }}
                  alt=""
                />
              ))}
            </Space>
          </Image.PreviewGroup>
        )}
      </section>

      {attrs && attrs.length > 0 ? (
        <section>
          <Typography.Title level={5}>{t("attributes")}</Typography.Title>
          <Descriptions bordered size="small" column={1} items={attrs} />
        </section>
      ) : null}

      {contact && contact.length > 0 ? (
        <section>
          <Typography.Title level={5}>{t("contactInfo")}</Typography.Title>
          <Descriptions bordered size="small" column={1} items={contact} />
        </section>
      ) : null}
    </Space>
  );
}
