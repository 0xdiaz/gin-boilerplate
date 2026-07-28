package auth

// EmailSender sends emails. Implement this interface to send password reset links in production.
// When injected into NewAuthService, ForgotPassword will send the token via email instead of returning it.
type EmailSender interface {
	SendPasswordResetEmail(to, resetToken string) error
}
