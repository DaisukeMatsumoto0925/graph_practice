import React, { useCallback, useState, useEffect, useMemo } from "react";
import {
  Dimmer,
  Loader,
  DropdownProps,
  Form,
  Message,
  Checkbox,
  Button,
  Icon
} from "semantic-ui-react";
import { useTaskFields } from "../../hooks/formHooks";

import {
  useCreateTaskMutation,
  Task
} from "../../src/generated/graphql";

const Tasks = () => {
  const {
    titleProps,
    notesProps
  } = useTaskFields();

  const [createTask] = useCreateTaskMutation({
    variables: {
      title: titleProps.value,
      note: notesProps.value,
    },
  });

  const handleButtonClick = useCallback(() => {
    createTask();
  }, [createTask]);

  return (
      <div>
        <h2>タスクの作成</h2>
        <Form>
          <Form.Field required={true}>
            <label>タスク名</label>
            <Form.Input
              placeholder="ピーマンを買いに行く"
              type="text"
              required={true}
            />
          </Form.Field>
          <Form.Field>
            <label>メモ</label>
            <Form.Input
              placeholder="駅前のOKストアがマジで安い"
              type="text"
            />
          </Form.Field>
        </Form>
        <Button
          icon={true}
          onClick={handleButtonClick}
          positive={true}
          >
            <Icon name="plus" /> 追加する
        </Button>
      </div>
  );
};
export default Tasks;
