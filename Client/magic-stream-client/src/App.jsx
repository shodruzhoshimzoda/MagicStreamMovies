import { useState } from 'react'

import './App.css'

import Home from './components/home/Home'
import Header from "./components/header/Header"
import Login from './components/login/Login'
import Register from './components/register/Register'
import {  Routes, Route, useNavigate} from 'react-router-dom'
import Recommended from './components/recommended/Recommended'
import Review from './components/review/Review'
import Layout from './components/Layout'
import RequiredAuth from './components/RequiredAuth'


function App() {
  const [count, setCount] = useState(0)

  const navigate = useNavigate();

  const updateMovieReview = (imdb_id) => { 
  navigate(`/review/${imdb_id}`);
}


  return (
    <>
    <Header/>
    <Routes path = "/" element = {<Layout/>}>
      <Route path='/' element={<Home updateMovieReview={updateMovieReview}/>}></Route>
      <Route path='/register' element={<Register />}></Route>
      <Route path='/login' element={<Login />}></Route>
      <Route element = {<RequiredAuth/> }>
          
              <Route path='/recommended' element={<Recommended />}></Route>
              <Route path='/review/:imdb_id' element={<Review />}></Route>



        </Route>
     

    </Routes>
    </>
  )
}

export default App
