async function shortenUrl() {
    const urlInput = document.getElementById('urlInput');
    const resultDiv = document.getElementById('result');
    const shortUrlLink = document.getElementById('shortUrl');
    
    if (!urlInput.value) {
        alert('Please enter a URL');
        return;
    }

    try {
        const response = await fetch('/shorten', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url: urlInput.value })
        });

        const data = await response.json();
        const shortUrl = data.short_url; // Extract the short_url from JSON
        
        shortUrlLink.href = shortUrl;
        shortUrlLink.textContent = shortUrl;
        resultDiv.classList.remove('hidden');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to shorten URL');
    }
}