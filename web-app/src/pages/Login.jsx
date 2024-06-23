import React from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../css/Login.css';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { LoginFormPage, ProFormText } from '@ant-design/pro-components';
import { theme, message } from 'antd';
import { login } from '../store/userSlice';
import { useDispatch } from 'react-redux';

const Login = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { token } = theme.useToken();

  const handleLogin = async (loginForm) => {
    try {
      const response = await axios.post('http://localhost:8080/login', loginForm);
      const { token, role, username} = response.data;

      if (role === 'broker') {
        dispatch(login({ username, token }));
        navigate('/broker');
      } else {
        message.warning('This account is not a broker, please try again');
      }
    } catch (err) {
      message.error('Login failed. Please try again.');
    }
  };

  return (
    <div style={{ backgroundColor: 'white', height: '100vh' }}>
      <LoginFormPage
        backgroundImageUrl="https://mdn.alipayobjects.com/huamei_gcee1x/afts/img/A*y0ZTS6WLwvgAAAAAAAAAAAAADml6AQ/fmt.webp"
        backgroundVideoUrl="https://gw.alipayobjects.com/v/huamei_gcee1x/afts/video/jXRBRK_VAwoAAAAAAAAAAAAAK4eUAQBr"
        title={<span style={{ color: 'white', fontWeight: 500 }}>Distributed Message Queue System</span>}
        containerStyle={{ backgroundColor: 'rgba(0, 0, 0,0.65)', backdropFilter: 'blur(4px)', paddingBottom: '50px' }}
        subTitle={<span style={{ color: 'white', fontWeight: 300 }}>Brain Stormtroopers</span>}
        onFinish={handleLogin}
      >
        <ProFormText
          name="username"
          fieldProps={{
            size: 'large',
            prefix: <UserOutlined style={{ color: token.colorText }} className={'prefixIcon'} />,
          }}
          placeholder="Please enter your username"
          rules={[{ required: true, message: 'Username cannot be empty!' }]}
        />
        <ProFormText.Password
          name="password"
          fieldProps={{
            size: 'large',
            prefix: <LockOutlined style={{ color: token.colorText }} className={'prefixIcon'} />,
          }}
          placeholder="Please enter your password"
          rules={[{ required: true, message: 'Password cannot be empty' }]}
        />
      </LoginFormPage>
    </div>
  );
};

export default Login;
