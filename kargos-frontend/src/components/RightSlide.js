import './css/RightSlideOpen.css'; // import the styles.css file

import { useState, useEffect } from "react";

/**
 * A function that generates a white rounded page sliding out from the right side.
 * @param {*}} props The props to use as argument.
 * @returns A button that opens the right slide page when clicked.
 */
export const RightSlidePage = (props) => {
    const [isOpen, setIsOpen] = useState(false);
    const { clickButton, content } = props;

    const openSpace = () => setIsOpen(true);
    const closeSpace = () => setIsOpen(false);
  
    useEffect(() => {
      const handleEsc = (event) => {
        if (event.key === 'Escape') {
          setIsOpen(false);
        }
      };
  
      window.addEventListener('keydown', handleEsc);
  
      return () => {
        window.removeEventListener('keydown', handleEsc);
      };
    }, []);
  
    return (
      <>
        <button className="open-button" onClick={openSpace}>
          {clickButton}
        </button>
        <div className={`overlay${isOpen ? ' open' : ''}`}></div>
        <div className={`space${isOpen ? ' open' : ''}`}>
          <button className="close-button" onClick={closeSpace}>
            X
          </button>
            <div>
            {content}
            </div>
        </div>
      </>
    );
}

