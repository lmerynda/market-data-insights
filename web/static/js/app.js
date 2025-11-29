// Basic JavaScript for Market Data Insights
document.addEventListener('DOMContentLoaded', function() {
    console.log('Market Data Insights loaded successfully!');
    
    // Add some interactivity to buttons
    const buttons = document.querySelectorAll('.btn');
    buttons.forEach(button => {
        button.addEventListener('mouseenter', function() {
            this.style.transform = 'scale(1.05)';
        });
        
        button.addEventListener('mouseleave', function() {
            this.style.transform = 'scale(1)';
        });
    });
});

// Function to check API health
async function checkAPIHealth() {
    try {
        const response = await fetch('/api/health');
        const data = await response.json();
        console.log('API Health:', data);
        return data;
    } catch (error) {
        console.error('Failed to check API health:', error);
        return null;
    }
}