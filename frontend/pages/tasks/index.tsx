import { NetworkStatus } from '@apollo/client';
import { useQuery } from "@apollo/client";
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
  Task,
  useTasksQuery,
} from "../../src/generated/graphql";

const Tasks = () => {
  const {
    titleProps,
    notesProps
  } = useTaskFields();

  const [createTask,{ data: createdData,loading: mutateLoading,called },] = useCreateTaskMutation({
    variables: {
      title: titleProps.value,
      note: notesProps.value,
    },
    onCompleted: () => {refetch()}
  });

  console.log(createdData)

  const handleButtonClick = useCallback(() => {
    createTask();
  }, [createTask]);

  // console.log(titleProps, notesProps)

  const {data, refetch, networkStatus, loading, error} = useTasksQuery()

  return (
      <>
        <h2>タスクの作成</h2>
        <Form>
          <Form.Field required={true}>
            <label>タスク名</label>
            <Form.Input
              placeholder="ピーマンを買いに行く"
              type="text"
              required={true}
              {...titleProps}
            />
          </Form.Field>
          <Form.Field>
            <label>メモ</label>
            <Form.Input
              placeholder="駅前のOKストアがマジで安い"
              type="text"
              {...notesProps}
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
        <div>
        <button onClick={() => refetch()}>Refetch!</button>
          {data?.tasks.map((task, i) => {
            return(
              <div key={i}>
                {task?.title}
                {" "}
                {task?.note}
                {" "}
                {task?.created_at}
              </div>
            )
          })}
        </div>
      </>
  );
};
export default Tasks;
