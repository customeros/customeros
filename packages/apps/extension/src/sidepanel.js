document.addEventListener('DOMContentLoaded', function() {
    const betaButton = document.getElementById('beta-button');
    
    betaButton.addEventListener('click', function() {
        const email = 'support@customeros.ai';
        const subject = 'LinkedIn helper request';
        const mailtoLink = `mailto:${email}?subject=${encodeURIComponent(subject)}`;
        
        window.open(mailtoLink, '_blank');
    });
});

//TODO: Implement side panel for LinkedIn pages