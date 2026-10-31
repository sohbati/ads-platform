import { Permission, type PermissionCode } from "@/lib/permissions/codes";

interface Rule {
  methods: readonly string[];
  pattern: RegExp;
  permission: PermissionCode;
}

const READ = ["GET"] as const;
const WRITE = ["POST", "PUT", "PATCH", "DELETE"] as const;
const ID = "\\d+";

// Paths are relative to /api/v1/panel. Anything not listed here is denied.
const RULES: Rule[] = [
  {
    methods: READ,
    pattern: new RegExp(`^staff(/${ID}(/(roles|groups|actions))?)?$`),
    permission: Permission.staffRead,
  },
  {
    methods: WRITE,
    pattern: new RegExp(`^staff(/${ID}(/password|/roles/${ID})?)?$`),
    permission: Permission.staffUpdate,
  },
  {
    methods: WRITE,
    pattern: new RegExp(`^staff-groups/${ID}/staff/${ID}$`),
    permission: Permission.staffUpdate,
  },
  {
    methods: READ,
    pattern: new RegExp(`^staff-groups(/${ID}(/staff)?)?$`),
    permission: Permission.staffGroupRead,
  },
  {
    methods: WRITE,
    pattern: new RegExp(`^staff-groups(/${ID})?$`),
    permission: Permission.staffGroupUpdate,
  },
  {
    methods: READ,
    pattern: new RegExp(`^roles(/${ID}(/actions)?)?$`),
    permission: Permission.roleRead,
  },
  {
    methods: WRITE,
    pattern: new RegExp(`^roles(/${ID}(/actions(/${ID})?)?)?$`),
    permission: Permission.roleUpdate,
  },
  {
    methods: READ,
    pattern: new RegExp(`^actions(/${ID})?$`),
    permission: Permission.actionRead,
  },
  {
    methods: WRITE,
    pattern: new RegExp(`^actions(/${ID})?$`),
    permission: Permission.actionUpdate,
  },
  {
    methods: READ,
    pattern: new RegExp(`^ads(/${ID})?$`),
    permission: Permission.advertisementRead,
  },
  {
    methods: ["POST"],
    pattern: new RegExp(`^ads/${ID}/approve$`),
    permission: Permission.advertisementApprove,
  },
  {
    methods: ["POST"],
    pattern: new RegExp(`^ads/${ID}/reject$`),
    permission: Permission.advertisementReject,
  },
  { methods: READ, pattern: /^settings$/, permission: Permission.settingsRead },
  {
    methods: ["PUT"],
    pattern: /^settings\/[a-z][a-z0-9_.]*$/,
    permission: Permission.settingsUpdate,
  },
];

/** The permission needed to call a panel path, or null if the path may not be proxied at all. */
export function requiredPermission(method: string, path: string): PermissionCode | null {
  const rule = RULES.find(
    (item) => item.methods.includes(method.toUpperCase()) && item.pattern.test(path),
  );
  return rule?.permission ?? null;
}

/** Writes that target one staff member, which only a super admin may make on another super admin. */
export function targetStaffId(method: string, path: string): number | null {
  if (method.toUpperCase() === "GET") {
    return null;
  }
  const match =
    /^staff\/(\d+)(\/|$)/.exec(path) ?? /^staff-groups\/\d+\/staff\/(\d+)$/.exec(path);
  return match ? Number(match[1]) : null;
}
