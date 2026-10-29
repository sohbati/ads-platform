import { expect, test } from "@playwright/test";

test("login reaches the dashboard", async ({ page }) => {
  await page.goto("/login");
  await page.getByLabel(/موبایل|Mobile/).fill("09120000000");
  await page.getByLabel(/رمز عبور|Password/).fill("admin");
  await page.getByRole("button", { name: /ورود|Sign in/ }).click();
  await expect(page).toHaveURL(/\/dashboard/);
});
