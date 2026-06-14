import { useState } from 'react'

import './App.css'

import Home from './components/home/Home'
import Header from "./components/header/Header"
import Login from './components/login/Login'
import Register from './components/register/Register'
import {  Routes, Route, useNavigate} from 'react-router-dom'


function App() {
  const [count, setCount] = useState(0)

  return (
    <>
    <Header/>
    <Routes>
      <Route path='/' element={<Home />}></Route>
      <Route path='/register' element={<Register />}></Route>
      <Route path='/login' element={<Login />}></Route>

    </Routes>
    </>
  )
}

export default App
