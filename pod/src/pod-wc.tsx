import React from 'react'
import ReactDOM from 'react-dom/client'
import reactToWebComponent from 'react-to-webcomponent'
import App from './App.tsx'
import './index.css'

const PodWebComponent = reactToWebComponent(App, React, ReactDOM)

customElements.define('pod-app', PodWebComponent)

export default PodWebComponent
