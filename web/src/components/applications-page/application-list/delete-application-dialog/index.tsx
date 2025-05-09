import { Application, selectById } from "~/modules/applications";
import {
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
  makeStyles,
} from "@material-ui/core";
import { FC, Fragment, memo, useCallback } from "react";
import { UI_TEXT_CANCEL, UI_TEXT_DELETE } from "~/constants/ui-text";
import {
  clearDeletingApp,
  deleteApplication,
} from "~/modules/delete-application";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";

import Alert from "@material-ui/lab/Alert";
import { DELETE_APPLICATION_SUCCESS } from "~/constants/toast-text";
import { Skeleton } from "@material-ui/lab";
import { addToast } from "~/modules/toasts";
import { red } from "@material-ui/core/colors";
import { shallowEqual } from "react-redux";
import { useStyles as useButtonStyles } from "~/styles/button";

const useStyles = makeStyles((theme) => ({
  applicationName: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
  },
  applicationLabels: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
  },
  description: {
    marginBottom: theme.spacing(2),
  },
  deleteButton: {
    color: theme.palette.getContrastText(red[400]),
    backgroundColor: red[800],
    "&:hover": {
      backgroundColor: red[800],
    },
  },
}));

const TITLE = "Delete Application";
const ALERT_TEXT = "Are you sure you want to delete the application?";

export interface DeleteApplicationDialogProps {
  onDeleted: () => void;
}

export const DeleteApplicationDialog: FC<DeleteApplicationDialogProps> = memo(
  function DeleteApplicationDialog({ onDeleted }) {
    const classes = useStyles();
    const buttonClasses = useButtonStyles();
    const dispatch = useAppDispatch();

    const [application, isDeleting] = useAppSelector<
      [Application.AsObject | undefined, boolean]
    >(
      (state) => [
        state.deleteApplication.applicationId
          ? selectById(
              state.applications,
              state.deleteApplication.applicationId
            )
          : undefined,
        state.deleteApplication.deleting,
      ],
      shallowEqual
    );

    const handleDelete = useCallback(() => {
      dispatch(deleteApplication()).then(() => {
        onDeleted();
        dispatch(
          addToast({ severity: "success", message: DELETE_APPLICATION_SUCCESS })
        );
      });
    }, [dispatch, onDeleted]);

    const handleCancel = useCallback(() => {
      dispatch(clearDeletingApp());
    }, [dispatch]);

    const renderLabels = useCallback(() => {
      if (!application?.labelsMap) return <Skeleton height={24} width={200} />;

      if (application.labelsMap.length === 0) return "-";

      return application.labelsMap.map(([key, value]) => (
        <Fragment key={key}>
          <span>{`${key}: ${value}`}</span>
          <br />
        </Fragment>
      ));
    }, [application?.labelsMap]);

    return (
      <Dialog
        open={Boolean(application)}
        onClose={(_event, reason) => {
          if (reason !== "backdropClick" || !isDeleting) {
            handleCancel();
          }
        }}
      >
        <DialogTitle>{TITLE}</DialogTitle>
        <DialogContent>
          <Alert severity="error" className={classes.description}>
            {ALERT_TEXT}
          </Alert>
          <Typography variant="caption">Name</Typography>
          <Typography variant="body1" className={classes.applicationName}>
            {application ? (
              application.name
            ) : (
              <Skeleton height={24} width={200} />
            )}
          </Typography>
          <Box height={24} />
          <Typography variant="caption">Labels</Typography>
          <Typography variant="body1" className={classes.applicationLabels}>
            {renderLabels()}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCancel} disabled={isDeleting}>
            {UI_TEXT_CANCEL}
          </Button>
          <Button
            variant="contained"
            color="primary"
            onClick={handleDelete}
            className={classes.deleteButton}
            disabled={isDeleting}
          >
            {UI_TEXT_DELETE}
            {isDeleting && (
              <CircularProgress size={24} className={buttonClasses.progress} />
            )}
          </Button>
        </DialogActions>
      </Dialog>
    );
  }
);
