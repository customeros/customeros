document.addEventListener("DOMContentLoaded", function () {
  const betaButton = document.getElementById("beta-button");

  if (!betaButton) return;

  betaButton.addEventListener("click", function () {
    const email = "support@customeros.ai";
    const subject = "LinkedIn helper request";
    const mailtoLink = `mailto:${email}?subject=${encodeURIComponent(subject)}`;

    window.open(mailtoLink, "_blank");
  });
});
