package impl

import (
	"context"
	"strconv"
	"testing"
	"time"

	"ads-platform/internal/business/otp/client"
	"ads-platform/internal/core/exception"
)

type fakeCache struct {
	values map[string]string
}

func newFakeCache() *fakeCache {
	return &fakeCache{values: map[string]string{}}
}

func (f *fakeCache) StoreOTP(_ context.Context, key string, otp string) error {
	f.values[key] = otp
	return nil
}

func (f *fakeCache) GetOTP(_ context.Context, key string) (string, error) {
	v, ok := f.values[key]
	if !ok {
		return "", &notFound{}
	}
	return v, nil
}

type notFound struct{}

func (*notFound) Error() string { return "not found" }

func TestSendOTPBlocksResendUntilWaitElapses(t *testing.T) {
	cache := newFakeCache()
	pub := &fakePublisher{}
	now := time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC)
	svc := NewOtpService(cache, pub, time.Minute).(*otpService)
	svc.now = func() time.Time { return now }

	first, err := svc.SendOTP(context.Background(), "+989121110001")
	if err != nil {
		t.Fatal(err)
	}
	if first.Message != "otp_sent" || first.ResendAfterSeconds != 60 {
		t.Fatalf("first=%+v", first)
	}
	if pub.calls != 1 {
		t.Fatalf("publish calls=%d, want 1", pub.calls)
	}
	otp := cache.values[otpCacheKey("+989121110001")]
	if len(otp) != 6 {
		t.Fatalf("otp=%q", otp)
	}

	_, err = svc.SendOTP(context.Background(), "+989121110001")
	appErr, ok := exception.AsAppError(err)
	if !ok || appErr.ErrorCode != "OTP_RESEND_WAIT" {
		t.Fatalf("err=%v", err)
	}
	if appErr.StatusCode != 429 || len(appErr.Params) != 1 || appErr.Params[0] != "60" {
		t.Fatalf("wait error: %+v", appErr)
	}
	if pub.calls != 1 {
		t.Fatalf("blocked send must not publish, calls=%d", pub.calls)
	}
	if cache.values[otpCacheKey("+989121110001")] != otp {
		t.Fatal("otp must not be replaced during wait")
	}

	now = now.Add(61 * time.Second)
	second, err := svc.SendOTP(context.Background(), "+989121110001")
	if err != nil {
		t.Fatal(err)
	}
	if second.Message != "otp_sent" {
		t.Fatalf("second=%+v", second)
	}
	if pub.calls != 2 {
		t.Fatalf("publish calls=%d, want 2", pub.calls)
	}
}

func TestSendOTPAllowsDifferentMobileDuringWait(t *testing.T) {
	cache := newFakeCache()
	pub := &fakePublisher{}
	svc := NewOtpService(cache, pub, time.Minute).(*otpService)
	svc.now = func() time.Time { return time.Unix(1_800_000_000, 0) }

	if _, err := svc.SendOTP(context.Background(), "+989121110001"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SendOTP(context.Background(), "+989121110002"); err != nil {
		t.Fatal(err)
	}
	if pub.calls != 2 {
		t.Fatalf("calls=%d", pub.calls)
	}
}

func TestCooldownRemainingParsesUnixUntil(t *testing.T) {
	cache := newFakeCache()
	svc := NewOtpService(cache, &fakePublisher{}, time.Minute).(*otpService)
	now := time.Unix(1_800_000_000, 0)
	svc.now = func() time.Time { return now }
	cache.values[cooldownCacheKey("+98912")] = strconv.FormatInt(now.Add(12*time.Second).Unix(), 10)
	if got := svc.cooldownRemaining(context.Background(), "+98912"); got != 12 {
		t.Fatalf("remaining=%d", got)
	}
}

func TestVerifyOTPRejectsAfterValidityWindow(t *testing.T) {
	cache := newFakeCache()
	now := time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC)
	svc := NewOtpService(cache, &fakePublisher{}, 2*time.Minute).(*otpService)
	svc.now = func() time.Time { return now }

	if _, err := svc.SendOTP(context.Background(), "+989121110001"); err != nil {
		t.Fatal(err)
	}
	otp := cache.values[otpCacheKey("+989121110001")]

	got, err := svc.VerifyOTP(context.Background(), "+989121110001", otp)
	if err != nil || got == nil || !got.Verified {
		t.Fatalf("valid window: err=%v got=%+v", err, got)
	}

	now = now.Add(2*time.Minute + time.Second)
	_, err = svc.VerifyOTP(context.Background(), "+989121110001", otp)
	appErr, ok := exception.AsAppError(err)
	if !ok || appErr.ErrorCode != "OTP_EXPIRED" {
		t.Fatalf("after timer err=%v", err)
	}
}

type fakePublisher struct {
	calls int
}

func (p *fakePublisher) PublishOtpEvent(context.Context, string, string) error {
	p.calls++
	return nil
}

var _ client.OtpEventPublisher = (*fakePublisher)(nil)
var _ client.OtpCacheClient = (*fakeCache)(nil)
