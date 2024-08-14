function handleEvent(eventType, selector, callback) {
    document.addEventListener(eventType, (event) => {
      const targetElement = event.target.closest(selector);
      if (targetElement) {
        callback(event, targetElement);
      }
    });
  }
  
  function createObjectFromHTML(html) {
    //   console.log(html)
    const wrapper = document.createElement("div");
    wrapper.innerHTML = html.trim();
  
    return createElementFromNode(wrapper.firstChild);
  }
  
  function createElementFromNode(node) {
    if (node.nodeType === Node.TEXT_NODE) {
      // Text node
      return node.textContent;
    }
  
    const tag = node.tagName.toLowerCase();
    const attrs = {};
    const children = [];
  
    // Collect attributes
    for (const { name, value } of node.attributes) {
      attrs[name] = value;
    }
  
    // Process child nodes
    for (const childNode of node.childNodes) {
      const childElement = createElementFromNode(childNode);
      children.push(childElement);
    }
  
    return { tag, attrs, children };
  }
  
  function render(element) {
    const domElement = document.createElement(element.tag);
  
    for (const [attr, value] of Object.entries(element.attrs)) {
      domElement.setAttribute(attr, value);
    }
  
    for (const child of element.children) {
      if (typeof child === "string") {
        const textNode = document.createTextNode(child);
        domElement.appendChild(textNode);
      } else {
        const childElement = render(child);
        domElement.appendChild(childElement);
      }
    }
  
    return domElement;
  }
  
  
  export { handleEvent, render, createObjectFromHTML };
  