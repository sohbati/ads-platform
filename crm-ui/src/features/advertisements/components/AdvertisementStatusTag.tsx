"use client";

import { Tag } from "antd";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { AdvertisementStatus, STATUS_LABEL } from "../constants/advertisement.constants";
import type { AdvertisementListItem } from "../types/advertisement.types";

const COLOR: Record<AdvertisementStatus, string> = {
  draft: "default",
  pending: "processing",
  active: "success",
  rejected: "error",
  expired: "default",
  deleted: "default",
};

export function AdvertisementStatusTag({
  ad,
}: {
  ad: Pick<AdvertisementListItem, "status" | "was_rejected">;
}) {
  const { t } = useI18n();
  // A rejected ad sits in draft until the customer fixes it; a fixed one is pending again.
  if (ad.was_rejected && ad.status === AdvertisementStatus.draft) {
    return <Tag color="error">{t("returnedToCustomer")}</Tag>;
  }
  return (
    <>
      <Tag color={COLOR[ad.status] ?? "default"}>
        {t(STATUS_LABEL[ad.status] ?? "status")}
      </Tag>
      {ad.was_rejected && ad.status === AdvertisementStatus.pending ? (
        <Tag color="warning">{t("resubmitted")}</Tag>
      ) : null}
    </>
  );
}
