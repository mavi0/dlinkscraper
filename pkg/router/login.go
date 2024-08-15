package router

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Login sends the username and password to the router for authentication.
// It expects to receive prompts for the username and password from the router,
// and writes the provided username and password to the router.
// If any error occurs during the process, it returns an error.
func (r *Router) Login(username, password string) error {
	
	if output, _, err := r.Expect("login: "); err != nil {
		return fmt.Errorf("failed to get username prompt from router %w", err)
	} else {
		logrus.WithField("output", output).Debugln("login field obtained from router")
	}

	if _, err := r.Write(fmt.Sprintf("%s\n", username)); err != nil {
		return fmt.Errorf("failed to write the username to the router: %w", err)
	}

	if output, _, err := r.Expect("Password: "); err != nil {
		return fmt.Errorf("failed to get password prompt from router: %w", err)
	} else {
		logrus.WithField("output", output).Debugln("password field obtained from router")
	}

	if _, err := r.Write(fmt.Sprintf("%s\n", password)); err != nil {
		return fmt.Errorf("failed to write the password to the router: %w", err)
	}
	return nil
}
