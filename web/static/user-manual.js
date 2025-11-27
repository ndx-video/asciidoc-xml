// User Manual Documentation Page Script
// This script handles XSLT transformation of the documentation

(function() {
    'use strict';

    // Get XML and XSLT content from script tags
    const xmlScript = document.getElementById('xml-data');
    const xsltScript = document.getElementById('xslt-data');
    const contentDiv = document.getElementById('content');

    if (!xmlScript || !xsltScript || !contentDiv) {
        console.error('Required elements not found');
        if (contentDiv) {
            contentDiv.innerHTML = '<p style="color: red;">Error: Required script elements not found</p>';
        }
        return;
    }

    let xmlContent, xsltContent;
    try {
        xmlContent = JSON.parse(xmlScript.textContent);
        xsltContent = JSON.parse(xsltScript.textContent);
    } catch (error) {
        console.error('Error parsing JSON data:', error);
        contentDiv.innerHTML = '<p style="color: red;">Error parsing data: ' + error.message + '</p>';
        return;
    }

    try {
        // Parse XML
        const parser = new DOMParser();
        const xmlDoc = parser.parseFromString(xmlContent, 'application/xml');
        
        // Check for parsing errors
        const parseError = xmlDoc.querySelector('parsererror');
        if (parseError) {
            console.error('XML parsing error:', parseError.textContent);
            contentDiv.innerHTML = '<p style="color: red;">Error parsing XML: ' + 
                parseError.textContent + '</p>';
            return;
        }

        // Parse XSLT
        const xsltDoc = parser.parseFromString(xsltContent, 'application/xml');
        
        // Check for parsing errors
        const xsltError = xsltDoc.querySelector('parsererror');
        if (xsltError) {
            console.error('XSLT parsing error:', xsltError.textContent);
            contentDiv.innerHTML = '<p style="color: red;">Error parsing XSLT: ' + 
                xsltError.textContent + '</p>';
            return;
        }

        // Transform
        const processor = new XSLTProcessor();
        processor.importStylesheet(xsltDoc);
        const result = processor.transformToFragment(xmlDoc, document);
        
        // Display
        contentDiv.appendChild(result);
    } catch (error) {
        console.error('Error during transformation:', error);
        contentDiv.innerHTML = '<p style="color: red;">Error: ' + error.message + '</p>';
    }
})();

