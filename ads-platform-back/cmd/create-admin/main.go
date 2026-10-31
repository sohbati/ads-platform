// Command create-admin creates (or promotes) a CRM super admin who can access every panel action.
//
//	go run ./cmd/create-admin -name "..." -mobile 0912... -national-id 0012345678
//
// The password is read from STAFF_ADMIN_PASSWORD, or prompted for without echo.
// Running it again for the same mobile promotes that staff member and resets their password.
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	staffrepoimpl "ads-platform/internal/business/iam/staff/repository/impl"
	"ads-platform/internal/business/iam/staff/service"
	staffserviceimpl "ads-platform/internal/business/iam/staff/service/impl"
	"ads-platform/internal/core/config"
	"ads-platform/internal/core/db"
	"ads-platform/internal/core/mobile"
)

func main() {
	name := flag.String("name", "", "display name (required for a new staff member)")
	mobileNumber := flag.String("mobile", "", "mobile number used to sign in (required)")
	nationalId := flag.String("national-id", "", "10-digit national ID (required for a new staff member)")
	flag.Parse()

	if strings.TrimSpace(*mobileNumber) == "" {
		flag.Usage()
		os.Exit(2)
	}

	password, err := readPassword()
	if err != nil {
		log.Fatalf("read password: %v", err)
	}

	cfg := config.Load()
	database, err := db.Connect(cfg.DatabaseURL, cfg.DatabaseType)
	if err != nil {
		log.Fatalf("database unavailable: %v", err)
	}

	staffService := staffserviceimpl.NewStaffService(
		staffrepoimpl.NewStaffRepository(database),
		mobile.NewNormalizer(cfg.DefaultCountryCode),
	)
	staff, err := staffService.UpsertSuperAdmin(context.Background(), service.StaffInput{
		Name:       *name,
		Mobile:     *mobileNumber,
		NationalId: *nationalId,
		Password:   password,
	})
	if err != nil {
		log.Fatalf("create admin: %v", err)
	}
	fmt.Printf("Super admin ready: id=%d name=%q mobile=%s\n", staff.ID, staff.Name, staff.Mobile)
}

func readPassword() (string, error) {
	if value := os.Getenv("STAFF_ADMIN_PASSWORD"); value != "" {
		return value, nil
	}
	first, err := prompt("Password: ")
	if err != nil {
		return "", err
	}
	second, err := prompt("Repeat password: ")
	if err != nil {
		return "", err
	}
	if first != second {
		return "", fmt.Errorf("passwords do not match")
	}
	return first, nil
}

var stdin = bufio.NewReader(os.Stdin)

func prompt(label string) (string, error) {
	fmt.Fprint(os.Stderr, label)
	setEcho(false)
	line, err := stdin.ReadString('\n')
	setEcho(true)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

// setEcho toggles terminal echo; it is a no-op when stdin is not a terminal.
func setEcho(on bool) {
	arg := "-echo"
	if on {
		arg = "echo"
	}
	cmd := exec.Command("stty", arg)
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}
