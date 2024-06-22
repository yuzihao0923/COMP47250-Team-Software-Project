import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../css/Login.css';
// import { Button, Checkbox, Form, Input } from 'antd';
import {
  AlipayOutlined,
  LockOutlined,
  MobileOutlined,
  TaobaoOutlined,
  UserOutlined,
  WeiboOutlined,
} from '@ant-design/icons';
import {
  LoginFormPage,
  ProConfigProvider,
  ProFormCaptcha,
  ProFormCheckbox,
  ProFormText,
} from '@ant-design/pro-components';
import { Button, Divider, Space, Tabs, message, theme } from 'antd';
import {login} from '../store/userSlice'
import { useDispatch } from 'react-redux';

/**
 * @type {React.CSSProperties}
 */
const mytyles = {
  // background: url('https://mdn.alipayobjects.com/huamei_gcee1x/afts/img/A*y0ZTS6WLwvgAAAAAAAAAAAAADml6AQ/fmt.webp') no-repeat center center fixed,
  // backgroundSize: 'cover',
  // height: '100%',
  // display: 'flex',
  // justifyContent: 'center',
  // alignItems: 'center'
};

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const dispatch = useDispatch()

  const { token } = theme.useToken();

  // const onFinish = (values) => {
  //   console.log('Success:', values);
  // };
  // const onFinishFailed = (errorInfo) => {
  //   console.log('Failed:', errorInfo);
  // };

  const handleLogin = async (loginForm) => {
    // console.log('Received values:', loginForm);
    const { username, password } = loginForm

    try {
      const response = await axios.post('http://localhost:8080/login', {
        username,
        password
      });
      console.log(response);

      const { token, role, username: user } = response.data;

      if (role === 'broker') {
        // localStorage.setItem(`${username}_token`, token); // Storing the JWT with username
        dispatch(login({user,token}))
        navigate('/broker');
      } else {
        setError('this user is not a broker, please try again');
        setUsername('');
        setPassword('');
      }
    } catch (err) {
      if (err.response && err.response.status === 401) {
        const errorMessage = err.response.data;
        setError(errorMessage);
        if (errorMessage.includes('username')) {
          setUsername('');
        }
        setPassword('');
      } else {
        setError('Login failed. Please try again.');
        setUsername('');
        setPassword('');
      }
    }
  };

  return (
    <div
      style={{
        backgroundColor: 'white',
        height: '100vh',
      }}
    >
      <LoginFormPage
        backgroundImageUrl="https://mdn.alipayobjects.com/huamei_gcee1x/afts/img/A*y0ZTS6WLwvgAAAAAAAAAAAAADml6AQ/fmt.webp"
        backgroundVideoUrl="https://gw.alipayobjects.com/v/huamei_gcee1x/afts/video/jXRBRK_VAwoAAAAAAAAAAAAAK4eUAQBr"
        title={<span style={{ color: 'white', fontWeight: 500 }}>Distributed Message Queue System</span>}
        containerStyle={{
          backgroundColor: 'rgba(0, 0, 0,0.65)',
          backdropFilter: 'blur(4px)',
          paddingBottom: '50px'
        }}
        subTitle={<span style={{ color: 'white', fontWeight: 300 }}>Brain Stormtroopers</span>}
        onFinish={handleLogin}
      >
        <ProFormText
          name="username"
          fieldProps={{
            size: 'large',
            prefix: (
              <UserOutlined
                style={{
                  color: token.colorText,
                }}
                className={'prefixIcon'}
              />
            ),
          }}
          placeholder={'Please enter your username'}
          rules={[
            {
              required: true,
              message: 'Username cannot be empty!',
            },
          ]}
        />
        <ProFormText.Password
          name="password"
          fieldProps={{
            size: 'large',
            prefix: (
              <LockOutlined
                style={{
                  color: token.colorText,
                }}
                className={'prefixIcon'}
              />
            ),
          }}
          placeholder={'Please enter your password'}
          rules={[
            {
              required: true,
              message: 'Password cannot be empty',
            },
          ]}
        />
        {/* <Button type='primary'>123</Button> */}
      </LoginFormPage>
    </div>
  )
};

export default Login;
