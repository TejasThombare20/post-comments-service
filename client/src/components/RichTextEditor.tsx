import React, { useState, useRef } from 'react';
import { Editor, EditorState, RichUtils, ContentState} from 'draft-js';
import { stateToHTML } from 'draft-js-export-html';
import { stateFromHTML } from 'draft-js-import-html';
import { Button } from '@/components/ui/button';
import 'draft-js/dist/Draft.css';

interface RichTextEditorProps {
  value?: string;
  onChange: (html: string) => void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
}

const RichTextEditor: React.FC<RichTextEditorProps> = ({
  value = '',
  onChange,
  placeholder = 'Write your comment...',
  className = '',
  disabled = false
}) => {
  const [editorState, setEditorState] = useState(() => {
    if (value) {
      try {
        // Check if the value is already HTML or plain text
        const isHTML = /<[a-z][\s\S]*>/i.test(value);
        if (isHTML) {
          const contentState = stateFromHTML(value);
          return EditorState.createWithContent(contentState);
        } else {
          // Convert plain text to HTML wrapped in span
          const htmlContent = `<span>${value}</span>`;
          const contentState = stateFromHTML(htmlContent);
          return EditorState.createWithContent(contentState);
        }
      } catch (error) {
        console.error('Error parsing content:', error);
        // Fallback to plain text
        const contentState = ContentState.createFromText(value);
        return EditorState.createWithContent(contentState);
      }
    }
    return EditorState.createEmpty();
  });

  const editorRef = useRef<Editor>(null);

  const handleEditorChange = (newEditorState: EditorState) => {
    setEditorState(newEditorState);
    
    // Convert to HTML and call onChange
    const contentState = newEditorState.getCurrentContent();
    const html = stateToHTML(contentState);
    onChange(html);
  };

  const handleKeyCommand = (command: string, editorState: EditorState) => {
    const newState = RichUtils.handleKeyCommand(editorState, command);
    if (newState) {
      handleEditorChange(newState);
      return 'handled';
    }
    return 'not-handled';
  };

  const toggleInlineStyle = (style: string) => {
    handleEditorChange(RichUtils.toggleInlineStyle(editorState, style));
  };

  const addLink = () => {
    const selection = editorState.getSelection();
    if (!selection.isCollapsed()) {
      const url = window.prompt('Enter URL:');
      if (url) {
        const contentState = editorState.getCurrentContent();
        const contentStateWithEntity = contentState.createEntity('LINK', 'MUTABLE', { url });
        const entityKey = contentStateWithEntity.getLastCreatedEntityKey();
        const newEditorState = EditorState.set(editorState, { currentContent: contentStateWithEntity });
        handleEditorChange(RichUtils.toggleLink(newEditorState, newEditorState.getSelection(), entityKey));
      }
    } else {
      alert('Please select text to add a link');
    }
  };

  const focusEditor = () => {
    if (editorRef.current) {
      editorRef.current.focus();
    }
  };

  // Custom style map for better styling
  const styleMap = {
    'BOLD': {
      fontWeight: 'bold',
    },
    'ITALIC': {
      fontStyle: 'italic',
    },
  };

  const currentStyle = editorState.getCurrentInlineStyle();

  return (
    <div className={`border border-gray-300 rounded-md ${className}`}>
      {/* Toolbar */}
      <div className="border-b border-gray-200 p-2 flex gap-1 bg-gray-50 rounded-t-md">
        <Button
          type="button"
          variant={currentStyle.has('BOLD') ? 'default' : 'ghost'}
          size="sm"
          onClick={() => toggleInlineStyle('BOLD')}
          disabled={disabled}
          className="h-8 w-8 p-0"
        >
          <strong>B</strong>
        </Button>
        <Button
          type="button"
          variant={currentStyle.has('ITALIC') ? 'default' : 'ghost'}
          size="sm"
          onClick={() => toggleInlineStyle('ITALIC')}
          disabled={disabled}
          className="h-8 w-8 p-0"
        >
          <em>I</em>
        </Button>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={addLink}
          disabled={disabled}
          className="h-8 px-2"
        >
          🔗
        </Button>
      </div>

      {/* Editor */}
      <div 
        className="p-3 min-h-[100px] cursor-text"
        onClick={focusEditor}
      >
        <Editor
          ref={editorRef}
          editorState={editorState}
          onChange={handleEditorChange}
          handleKeyCommand={handleKeyCommand}
          placeholder={placeholder}
          customStyleMap={styleMap}
          readOnly={disabled}
        />
      </div>
    </div>
  );
};

export default RichTextEditor; 