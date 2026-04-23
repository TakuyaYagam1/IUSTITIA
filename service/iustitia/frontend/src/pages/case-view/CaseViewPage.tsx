import { ROUTES } from '@/app/router/routes';
import { CaseView } from '@entities/case';
import { useAuthStore } from '@entities/user';
import { useCase, useCaseComplaints, useCaseOpinion } from '@features/case-view';
import { Button } from '@shared/ui';
import { ChevronLeft, FileText, Gavel } from 'lucide-react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import styles from './CaseViewPage.module.css';

export const CaseViewPage = (): JSX.Element => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const session = useAuthStore((s) => s.session);
  const role = session?.role;
  const userId = session?.userId;
  const canSeeComplaints = role === 'judge' || role === 'prosecutor';
  const canSeeOpinion = role === 'judge' || role === 'prosecutor';

  const caseQuery = useCase(id);
  const complaintsQuery = useCaseComplaints(canSeeComplaints ? id : undefined);
  useCaseOpinion(canSeeOpinion ? id : undefined, caseQuery.data?.status);

  const caseData = caseQuery.data;

  return (
    <div className={styles.wrap}>
      <Link to={ROUTES.cases} className={styles.back}>
        <ChevronLeft size={14} /> К картотеке
      </Link>

      {caseQuery.isLoading && <div className={styles.message}>Получение дела…</div>}
      {caseQuery.isError && (
        <div className={`${styles.message} ${styles.error}`}>
          Ошибка: {caseQuery.error instanceof Error ? caseQuery.error.message : 'неизвестно'}
        </div>
      )}
      {caseData && (
        <CaseView
          caseItem={caseData}
          {...(canSeeComplaints && complaintsQuery.data
            ? { complaints: complaintsQuery.data }
            : {})}
        />
      )}

      {caseData &&
        id &&
        role === 'prosecutor' &&
        caseData.status === 'assigned' &&
        caseData.assigned_prosecutor_id === userId && (
          <div className={styles.actions}>
            <Button variant="primary" onClick={() => navigate(ROUTES.caseOpinion(id))}>
              <FileText size={14} /> Подать заключение
            </Button>
          </div>
        )}

      {caseData && id && role === 'judge' && caseData.status === 'hearing' && (
        <div className={styles.actions}>
          <Button variant="primary" onClick={() => navigate(ROUTES.caseVerdict(id))}>
            <Gavel size={14} /> Вынести приговор
          </Button>
        </div>
      )}
    </div>
  );
};
