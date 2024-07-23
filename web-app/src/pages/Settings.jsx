import React from 'react';
import { Button, Form, InputNumber } from 'antd';

export default function Settings() {

  const onFinish = (values) => {
    console.log('Success:', values);
  };
  const onFinishFailed = (errorInfo) => {
    console.log('Failed:', errorInfo);
  };
  const onChange = (value) => {
    console.log('changed', value);
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
          <InputNumber min={0} defaultValue={0} onChange={onChange} />;
        </Form.Item>
      </Form>
    </div>
  )
}
