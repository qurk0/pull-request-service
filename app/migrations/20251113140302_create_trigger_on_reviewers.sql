-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_pr_reviewers_limit()
RETURNS TRIGGER AS $$
DECLARE
    reviewers_count INTEGER;
BEGIN
    SELECT COUNT(*)
    INTO reviewers_count
    FROM pull_requests_reviewers
    WHERE pull_request_id = NEW.pull_request_id
        AND (TG_OP <> 'UPDATE' OR NOT (pull_request_id = OLD.pull_request_id 
                                        AND reviewer_id = OLD.reviewer_id))
    IF reviewers_count >= 2 THEN
        RAISE EXCEPTION
            'cannot assign more than 2 reviewers to pull_request_id=%', NEW.pull_request_id
            USING ERRCODE = 'check_violation';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_limit_pr_reviewers
BEFORE INSERT OR UPDATE ON pull_requests_reviewers
FOR EACH ROW
EXECUTE FUNCTION check_pr_reviewers_limit();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_limit_pr_reviewers ON pull_requests_reviewers;
DROP FUNCTION IF EXISTS check_pr_reviewers_limit();
-- +goose StatementEnd
