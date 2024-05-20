package router

import (
	"fmt"
)

func (r *Router) Login(username, password string) error {
	if _, _, err := r.Expect("login: "); err != nil {
		return fmt.Errorf("failed to get username prompt from router %w", err)
	}
	if _, err := r.Write(fmt.Sprintf("%s\n", username)); err != nil {
		return fmt.Errorf("failed to write the username to the router: %w", err)
	}
	if _, _, err := r.Expect("Password: "); err != nil {
		return fmt.Errorf("failed to get password prompt from router: %w", err)
	}
	if _, err := r.Write(fmt.Sprintf("%s\n", password)); err != nil {
		return fmt.Errorf("failed to write the password to the router: %w", err)
	}
	return nil
}
