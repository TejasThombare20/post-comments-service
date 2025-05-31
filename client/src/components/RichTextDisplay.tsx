import React from 'react';
import DOMPurify from 'dompurify';

interface RichTextDisplayProps {
  content: string;
  className?: string;
}

const RichTextDisplay: React.FC<RichTextDisplayProps> = ({ content, className = '' }) => {
  // Function to convert plain text to HTML if needed
  const processContent = (rawContent: string): string => {
    if (!rawContent) return '';
    
    // Check if content is already HTML
    const isHTML = /<[a-z][\s\S]*>/i.test(rawContent);
    
    if (isHTML) {
      // Content is already HTML, sanitize it
      return DOMPurify.sanitize(rawContent, {
        ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'b', 'i', 'a', 'span', 'div'],
        ALLOWED_ATTR: ['href', 'target', 'rel']
      });
    } else {
      // Content is plain text, wrap in span and preserve line breaks
      const textWithBreaks = rawContent.replace(/\n/g, '<br>');
      const htmlContent = `<span>${textWithBreaks}</span>`;
      return DOMPurify.sanitize(htmlContent, {
        ALLOWED_TAGS: ['span', 'br'],
        ALLOWED_ATTR: []
      });
    }
  };

  const sanitizedContent = processContent(content);

  return (
    <div 
      className={`rich-text-content ${className}`}
      dangerouslySetInnerHTML={{ __html: sanitizedContent }}
      style={{
        wordBreak: 'break-word',
        lineHeight: '1.5'
      }}
    />
  );
};

export default RichTextDisplay; 