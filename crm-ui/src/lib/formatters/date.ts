import { formatDigits } from "./digits";

// The backend sends moments in UTC; they show in the viewer's own time zone, in the Jalali
// calendar for fa and the Gregorian one for en. Dates come from client-side queries, so they
// render in the browser and use its zone.

const formats = new Map<string, Intl.DateTimeFormat>();

function formatterFor(locale: "fa" | "en", timeZone: string | undefined): Intl.DateTimeFormat {
  const key = `${locale}|${timeZone ?? ""}`;
  let format = formats.get(key);
  if (!format) {
    // Latin digits here; formatDigits turns them into Persian digits for fa.
    format = new Intl.DateTimeFormat(locale === "fa" ? "en-US-u-ca-persian-nu-latn" : "en-GB", {
      timeZone,
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      hourCycle: "h23",
    });
    formats.set(key, format);
  }
  return format;
}

interface DateParts {
  year: string;
  month: string;
  day: string;
  hour: string;
  minute: string;
}

const PART_TYPES = new Set(["year", "month", "day", "hour", "minute"]);

function parse(value: string | Date | null | undefined): Date | null {
  if (!value) {
    return null;
  }
  const date = typeof value === "string" ? new Date(value) : value;
  return Number.isNaN(date.getTime()) ? null : date;
}

function partsOf(date: Date, locale: "fa" | "en", timeZone: string | undefined): DateParts {
  const parts: Partial<DateParts> = {};
  for (const part of formatterFor(locale, timeZone).formatToParts(date)) {
    if (PART_TYPES.has(part.type)) {
      parts[part.type as keyof DateParts] = part.value;
    }
  }
  return parts as DateParts;
}

/**
 * A date as yyyy/mm/dd in the Jalali calendar (fa) or dd/mm/yyyy (en).
 * `timeZone` defaults to the viewer's own.
 */
export function formatDate(
  value: string | Date | null | undefined,
  locale: "fa" | "en",
  timeZone?: string,
): string {
  const date = parse(value);
  if (!date) {
    return "—";
  }
  const { year, month, day } = partsOf(date, locale, timeZone);
  return locale === "fa"
    ? formatDigits(`${year}/${month}/${day}`, "fa")
    : `${day}/${month}/${year}`;
}

/** formatDate followed by the time as HH:mm. */
export function formatDateTime(
  value: string | Date | null | undefined,
  locale: "fa" | "en",
  timeZone?: string,
): string {
  const date = parse(value);
  if (!date) {
    return "—";
  }
  const { hour, minute } = partsOf(date, locale, timeZone);
  const time = `${hour}:${minute}`;
  return `${formatDate(date, locale, timeZone)} ${locale === "fa" ? formatDigits(time, "fa") : time}`;
}
