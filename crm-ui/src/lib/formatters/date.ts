import { formatDigits } from "./digits";

const GREGORIAN_EPOCH = 1721425.5;
const JALALI_EPOCH = 1948320.5;

function div(a: number, b: number): number {
  return Math.trunc(a / b);
}

function gregorianToJd(year: number, month: number, day: number): number {
  return (
    GREGORIAN_EPOCH -
    1 +
    365 * (year - 1) +
    div(year - 1, 4) -
    div(year - 1, 100) +
    div(year - 1, 400) +
    div(367 * month - 362, 12) +
    (month <= 2 ? 0 : leapGregorian(year) ? -1 : -2) +
    day
  );
}

function leapGregorian(year: number): boolean {
  return year % 4 === 0 && (year % 100 !== 0 || year % 400 === 0);
}

function jdToJalali(jd: number): [number, number, number] {
  const depoch = Math.floor(jd) + 0.5 - JALALI_EPOCH;
  const cycle = div(depoch, 1029983);
  const cyear = depoch % 1029983;
  let ycycle: number;
  if (cyear === 1029982) {
    ycycle = 2820;
  } else {
    const aux1 = div(cyear, 366);
    const aux2 = cyear % 366;
    ycycle = div(2134 * aux1 + 2816 * aux2 + 2816, 1028522) + aux1 + 1;
  }
  let year = ycycle + 2820 * cycle + 474;
  if (year <= 0) {
    year -= 1;
  }
  const yday = Math.floor(jd) + 0.5 - jalaliToJd(year, 1, 1) + 1;
  const month = yday <= 186 ? Math.ceil(yday / 31) : Math.ceil((yday - 6) / 30);
  const day = Math.floor(jd) + 0.5 - jalaliToJd(year, month, 1) + 1;
  return [year, month, day];
}

function jalaliToJd(year: number, month: number, day: number): number {
  const epbase = year - (year >= 0 ? 474 : 473);
  const epyear = 474 + (epbase % 2820);
  return (
    day +
    (month <= 7 ? (month - 1) * 31 : (month - 1) * 30 + 6) +
    div(epyear * 682 - 110, 2816) +
    (epyear - 1) * 365 +
    div(epbase, 2820) * 1029983 +
    (JALALI_EPOCH - 1)
  );
}

function pad(value: number): string {
  return String(value).padStart(2, "0");
}

export function formatDate(
  value: string | Date | null | undefined,
  locale: "fa" | "en",
): string {
  if (!value) {
    return "—";
  }
  const date = typeof value === "string" ? new Date(value) : value;
  if (Number.isNaN(date.getTime())) {
    return "—";
  }
  if (locale === "en") {
    return new Intl.DateTimeFormat("en-GB", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
    }).format(date);
  }
  const jd = gregorianToJd(
    date.getFullYear(),
    date.getMonth() + 1,
    date.getDate(),
  );
  const [year, month, day] = jdToJalali(jd);
  return formatDigits(`${year}/${pad(month)}/${pad(day)}`, "fa");
}

export function formatDateTime(
  value: string | Date | null | undefined,
  locale: "fa" | "en",
): string {
  if (!value) {
    return "—";
  }
  const date = typeof value === "string" ? new Date(value) : value;
  if (Number.isNaN(date.getTime())) {
    return "—";
  }
  const time = `${pad(date.getHours())}:${pad(date.getMinutes())}`;
  return `${formatDate(date, locale)} ${locale === "fa" ? formatDigits(time, "fa") : time}`;
}
