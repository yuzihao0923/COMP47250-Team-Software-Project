import React from 'react';
import { Button, Form, InputNumber } from 'antd';
import { changeWarningNumber } from '../store/settings'
import { useDispatch, useSelector } from 'react-redux';

export default function Settings() {

  const dispatch = useDispatch()
  const warningNumber = useSelector((state) => state.settings.warningNumber);

  const onFinish = (values) => {
    console.log('Success:', values);
  };
  const onFinishFailed = (errorInfo) => {
    console.log('Failed:', errorInfo);
  };
  const onChange = (value) => {
    console.log('changed', value);
    dispatch(changeWarningNumber({ warningNumber: value }))
  };

  return (
    <div>
      <Form
        name="settingsForm"
        labelCol={{
          span: 8,
        }}
        wrapperCol={{
          span: 16,
        }}
        style={{
          maxWidth: 600,
        }}
        initialValues={{
          remember: true,
        }}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        autoComplete="off"
      >
        <Form.Item
          label="warning number"
          name="warningNumber"
        >
          <InputNumber min={0} value={warningNumber} onChange={onChange} />;
        </Form.Item>
      </Form>
    </div>
  )
}
