(function () {
  const digitMap = {
    "۰": "0", "۱": "1", "۲": "2", "۳": "3", "۴": "4",
    "۵": "5", "۶": "6", "۷": "7", "۸": "8", "۹": "9",
    "٠": "0", "١": "1", "٢": "2", "٣": "3", "٤": "4",
    "٥": "5", "٦": "6", "٧": "7", "٨": "8", "٩": "9",
  };

  function normalizeCountryCode(code) {
    const trimmed = (code || "+98").trim();
    if (!trimmed) return "+98";
    const digits = trimmed.replace(/^\+/, "").replace(/\D/g, "");
    return "+" + digits;
  }

  function convertDigits(input) {
    return String(input).replace(/[۰-۹٠-٩]/g, function (ch) {
      return digitMap[ch] || ch;
    });
  }

  function stripSeparators(input) {
    return String(input).replace(/[\s\-()]/g, "");
  }

  function digitsOnly(input) {
    return String(input).replace(/\D/g, "");
  }

  function validateMobile(normalized, ccDigits) {
    if (!normalized.startsWith("+")) {
      return false;
    }

    const digits = digitsOnly(normalized.slice(1));
    if (!digits) {
      return false;
    }

    if (ccDigits === "98") {
      return digits.length === 12 && digits.startsWith("98") && digits[2] === "9";
    }

    return digits.startsWith(ccDigits) && digits.length > ccDigits.length;
  }

  function normalizeMobile(input, defaultCountryCode) {
    const defaultCC = normalizeCountryCode(defaultCountryCode);
    const ccDigits = defaultCC.slice(1);

    let value = stripSeparators(convertDigits(String(input || "").trim()));
    if (!value) {
      return { ok: false, error: "MOBILE_EMPTY" };
    }

    if (value.startsWith("00")) {
      value = "+" + value.slice(2);
    }

    let normalized = "";
    if (value.startsWith("+")) {
      normalized = "+" + digitsOnly(value.slice(1));
    } else if (value.startsWith("0")) {
      normalized = defaultCC + value.slice(1);
    } else if (value.startsWith(ccDigits)) {
      normalized = "+" + value;
    } else {
      normalized = defaultCC + value;
    }

    if (!validateMobile(normalized, ccDigits)) {
      return { ok: false, error: "INVALID_MOBILE" };
    }

    return { ok: true, mobile: normalized };
  }

  window.normalizeMobile = normalizeMobile;
})();
