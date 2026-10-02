package impl

import (
	"context"
	"testing"

	customermodel "ads-platform/internal/business/customer/model"
	"ads-platform/internal/business/customerprofile/model"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type fakeProfileRepo struct {
	byCustomer map[int64]*model.CustomerProfile
}

func (f *fakeProfileRepo) GetByCustomerID(_ context.Context, customerID int64) (*model.CustomerProfile, error) {
	p, ok := f.byCustomer[customerID]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *p
	cp.LocationSlugs = append([]string(nil), p.LocationSlugs...)
	return &cp, nil
}

func (f *fakeProfileRepo) Upsert(_ context.Context, profile *model.CustomerProfile) error {
	cp := *profile
	cp.LocationSlugs = append([]string(nil), profile.LocationSlugs...)
	f.byCustomer[profile.CustomerID] = &cp
	return nil
}

type fakeCustomers struct {
	ids map[int64]struct{}
}

func (f *fakeCustomers) GetCustomerByID(_ context.Context, id int64) (*customermodel.Customer, error) {
	if _, ok := f.ids[id]; !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return &customermodel.Customer{ID: id}, nil
}

func newService() *profileService {
	return &profileService{
		profiles:  &fakeProfileRepo{byCustomer: map[int64]*model.CustomerProfile{}},
		customers: &fakeCustomers{ids: map[int64]struct{}{1: {}}},
	}
}

func TestGetMissingProfileReturnsEmptySlugs(t *testing.T) {
	svc := newService()
	got, err := svc.Get(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if got.CustomerID != 1 || len(got.LocationSlugs) != 0 {
		t.Fatalf("got %+v", got)
	}
}

func TestGetUnknownCustomer(t *testing.T) {
	svc := newService()
	_, err := svc.Get(context.Background(), 99)
	assertCode(t, err, "CUSTOMER_NOT_FOUND")
}

func TestPutNormalizesAndDedupes(t *testing.T) {
	svc := newService()
	got, err := svc.Put(context.Background(), 1, []string{" Tehran ", "karaj", "tehran", ""})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.LocationSlugs) != 2 || got.LocationSlugs[0] != "tehran" || got.LocationSlugs[1] != "karaj" {
		t.Fatalf("slugs=%v", got.LocationSlugs)
	}
}

func TestPutRejectsBadSlug(t *testing.T) {
	svc := newService()
	_, err := svc.Put(context.Background(), 1, []string{"tehran!"})
	assertCode(t, err, "PROFILE_INVALID_LOCATION")
}

func TestPutRejectsTooMany(t *testing.T) {
	svc := newService()
	in := make([]string, maxLocationSlugs+1)
	for i := range in {
		in[i] = "city-" + itoa(i)
	}
	_, err := svc.Put(context.Background(), 1, in)
	assertCode(t, err, "PROFILE_TOO_MANY_LOCATIONS")
}

func TestPutInvalidCustomerID(t *testing.T) {
	svc := newService()
	_, err := svc.Put(context.Background(), 0, nil)
	assertCode(t, err, "PROFILE_INVALID_CUSTOMER")
}

func assertCode(t *testing.T, err error, code string) {
	t.Helper()
	app, ok := exception.AsAppError(err)
	if !ok {
		t.Fatalf("err=%v", err)
	}
	if app.ErrorCode != code {
		t.Fatalf("code=%s want %s", app.ErrorCode, code)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [12]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
